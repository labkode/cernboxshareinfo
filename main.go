package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/knq/dburl"
	"github.com/labkode/cernboxshareinfo/models"
	"github.com/labkode/cernboxshareinfo/proto"
	"github.com/syndtr/goleveldb/leveldb"
)

var sqlFile string
var phonebookBinary string
var phonebookdb string
var mysqlhost string
var mysqlusername string
var mysqlpassword string
var mysqldatabase string
var mysqlport int

func main() {
	flag.StringVar(&sqlFile, "sqlfile", "shares.sql", "file containing the mysql oc_shares dump")
	flag.StringVar(&phonebookBinary, "phonebookbinary", "/usr/bin/phonebook", "the phonebook binary")
	flag.StringVar(&phonebookdb, "phonebookdb", "phonebook.db", "phonebook cache db file")
	flag.StringVar(&mysqlhost, "mysqlhost", "localhost", "mysql hostname")
	flag.StringVar(&mysqlusername, "mysqlusername", "admin", "mysql username")
	flag.StringVar(&mysqlpassword, "mysqlpassword", "admin", "mysql password")
	flag.StringVar(&mysqldatabase, "mysqldatabase", "mydb", "mysql database")
	flag.IntVar(&mysqlport, "mysqlport", 3306, "mysql port")
	flag.Parse()

	fmt.Printf("Loading shares from: %s\n", sqlFile)
	fmt.Printf("PhoneBook Application: %s\n", phonebookBinary)

	handle, err := os.Open(sqlFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer handle.Close()

	r := csv.NewReader(handle)
	r.Comma = '\t'

	shareInfos := []*proto.ShareInfo{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(os.Stderr, err)
		} else {

			/*
			*
			* These are the attributes from the MySQL oc_share dump:
			*
			* id: unique identifier of the share
			* share_type: the type of the share: 0 is internal, 1 is group and 3 is by link
			* share_with: the receiver of the share
			* uid_owner: the username that started the share
			* parent: the parent file of the shared one
			* item_type: the type of resource being shared: either file or folder
			* item source: inode of file being shared
			* item_target: name of the file when appears in the sharee filesystem
			* file_source: the inode of the file being shared
			* file_target: the name of the file when appears in the sharee filesystem
			* permissions: the permissions for the sharee
			* stime: the time the share was created
			* accepted: the share was acepted: only valid for federated shares
			* expiration: when the share will expire: only valid for link shares
			* token: the token to access the share: only valid for link shares
			* mail_send: whether the sharee has been informed by email of the new share
			 */

			id, _ := strconv.Atoi(record[0])
			shareType, _ := strconv.Atoi(record[1])
			parent, _ := strconv.Atoi(record[4])
			itemSource, _ := strconv.Atoi(record[6])
			fileSource, _ := strconv.Atoi(record[8])
			stime, _ := strconv.Atoi(record[11])
			accepted, _ := strconv.Atoi(record[12])
			mailSend, _ := strconv.Atoi(record[15])

			shareInfo := &proto.ShareInfo{}
			shareInfo.Id = int64(id)
			shareInfo.ShareType = int64(shareType)
			shareInfo.ShareWith = record[2]
			shareInfo.UidOwner = record[3]
			shareInfo.Parent = int64(parent)
			shareInfo.ItemType = record[5]
			shareInfo.ItemSource = int64(itemSource)
			shareInfo.ItemTarget = record[7]
			shareInfo.FileSource = int64(fileSource)
			shareInfo.FileTarget = record[9]
			shareInfo.Permissions = record[10]
			shareInfo.Stime = int64(stime)
			shareInfo.Accepted = int64(accepted)
			shareInfo.Expiration = record[13]
			shareInfo.Token = record[14]
			shareInfo.MailSend = int64(mailSend)

			// we are interested in user shares
			if shareInfo.ShareType == 0 {
				shareInfos = append(shareInfos, shareInfo)
			}
		}
	}

	userMap := map[string]*proto.PersonInfo{}
	db, err := leveldb.OpenFile(phonebookdb, nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println("phonebook calls won't be cached")
	} else {
		fmt.Println("Loading phonebook db into memory...")
		total := 0
		iter := db.NewIterator(nil, nil)
		for iter.Next() {
			key := iter.Key()
			value := iter.Value() // value is json
			personInfo := &proto.PersonInfo{}
			err := json.Unmarshal(value, personInfo)
			if err != nil {
				fmt.Println(err)
			} else {
				userMap[string(key)] = personInfo
				total++
			}
		}
		iter.Release()
		fmt.Printf("Loaded %d entries from phonebook cached db\n", total)
	}

	total := len(shareInfos)
	for i, shareInfo := range shareInfos {
		fmt.Printf("Resolving users (%d)/%d\n", i, total)
		if _, ok := userMap[shareInfo.UidOwner]; !ok {
			personInfo, err := getPersonInfo(shareInfo.UidOwner)
			if err != nil {
				log.Println(err)
				db.Delete([]byte(shareInfo.UidOwner), nil)
			} else {
				userMap[shareInfo.UidOwner] = personInfo
				if db != nil {
					json, err := json.Marshal(personInfo)
					if err != nil {
						fmt.Println(err)
					} else {
						db.Put([]byte(shareInfo.UidOwner), json, nil)
					}
				}
			}
		}
		if _, ok := userMap[shareInfo.ShareWith]; !ok {
			personInfo, err := getPersonInfo(shareInfo.ShareWith)
			if err != nil {
				log.Println(err)
				db.Delete([]byte(shareInfo.UidOwner), nil)
			} else {
				userMap[shareInfo.ShareWith] = personInfo
				if db != nil {
					json, err := json.Marshal(personInfo)
					if err != nil {
						fmt.Println(err)
					} else {
						db.Put([]byte(shareInfo.ShareWith), json, nil)
					}
				}
			}
		}
	}

	fmt.Printf("Found %d persons with valid information in phonebook\n", len(userMap))

	flatInfos := []*proto.FlatInfo{}
	// create flat infos with owner and sharee together with share info
	for _, shareInfo := range shareInfos {
		_, ownerFound := userMap[shareInfo.UidOwner]
		_, shareeFound := userMap[shareInfo.ShareWith]
		if ownerFound && shareeFound {
			flatInfo := &proto.FlatInfo{}
			flatInfo.ShareInfo = shareInfo
			flatInfo.OwnerInfo = userMap[shareInfo.UidOwner]
			flatInfo.ShareeInfo = userMap[shareInfo.ShareWith]
			flatInfos = append(flatInfos, flatInfo)
		}
	}

	// forward data to sql
	sqldb, err := dburl.Open(fmt.Sprintf("mysql://%s:%s@%s:%d/%s", mysqlusername, mysqlpassword, mysqlhost, mysqlport, mysqldatabase))
	if err != nil {
		log.Fatal(err)
	}

	createCompanies(sqldb, flatInfos)
}

func createCompanies(db *sql.DB, flatInfos []*proto.FlatInfo) {
	companies := map[string]*proto.FlatInfo{}
	for _, flatInfo := range flatInfos {
		if _, ok := companies[flatInfo.GetOwnerInfo().Company]; !ok {
			companies[flatInfo.GetOwnerInfo().Company] = flatInfo
		}
		if _, ok := companies[flatInfo.GetShareeInfo().Company]; !ok {
			companies[flatInfo.GetShareeInfo().Company] = flatInfo
		}
	}
	for pk := range companies {
		company := &models.DimensionCompany{}
		company.Company = pk
		if err := company.Save(db); err != nil {
			log.Fatal(err)
		}
	}
}
func getPersonInfo(username string) (*proto.PersonInfo, error) {
	cmd := exec.Command(phonebookBinary, "--login", username, "-t", "login", "-t", "uid", "-t", "department", "-t", "organization", "-t", "company", "-t", "office")
	stdout, _, err := executeCMD(cmd)
	if err != nil {
		return nil, err
	}
	if stdout == "" {
		return nil, errors.New("user " + username + " does not exist anymore")
	}
	personInfo := &proto.PersonInfo{}
	tokens := strings.Split(stdout, ";")
	uid, _ := strconv.Atoi(tokens[1])
	personInfo.Login = tokens[0]
	personInfo.Uid = int64(uid)
	personInfo.Department = tokens[2]
	personInfo.Organization = tokens[3]
	personInfo.Company = tokens[4]
	personInfo.Office = tokens[5]
	return personInfo, nil
}

func executeCMD(cmd *exec.Cmd) (string, string, error) {
	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf
	err := cmd.Run()
	return outBuf.String(), errBuf.String(), err
}
