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
	"time"

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
var testcases int
var fillstarmodel bool
var showdata bool
var outputasjson bool
var version bool

// Build information obtained with the help of -ldflags
var (
	appName       string = "cernboxshareinfo"
	buildDate     string // date -u
	gitTag        string // git describe --exact-match HEAD
	gitNearestTag string // git describe --abbrev=0 --tags HEAD
	gitCommit     string // git rev-parse HEAD
)

func main() {
	flag.StringVar(&sqlFile, "sqlfile", "shares.sql", "file containing the mysql oc_shares dump")
	flag.StringVar(&phonebookBinary, "phonebookbinary", "/usr/bin/phonebook", "the phonebook binary")
	flag.StringVar(&phonebookdb, "phonebookdb", "phonebook.db", "phonebook cache db file")
	flag.StringVar(&mysqlhost, "mysqlhost", "localhost", "mysql hostname")
	flag.StringVar(&mysqlusername, "mysqlusername", "admin", "mysql username")
	flag.StringVar(&mysqlpassword, "mysqlpassword", "admin", "mysql password")
	flag.StringVar(&mysqldatabase, "mysqldatabase", "mydb", "mysql database")
	flag.IntVar(&mysqlport, "mysqlport", 3306, "mysql port")
	flag.IntVar(&testcases, "testcases", -1, "number of test cases")
	flag.BoolVar(&fillstarmodel, "fillstarmodel", false, "fill the existing sql star model on the configured mysql database")
	flag.BoolVar(&showdata, "showdata", true, "if enabled outputs the correlated data (one line per record)")
	flag.BoolVar(&outputasjson, "outputasjson", false, "if enabled outputs the data in JSON format instead of proto format")
	flag.BoolVar(&version, "version", false, "show version")
	flag.Parse()

	if version {
		// if gitTag is not empty we are on release build
		if gitTag != "" {
			fmt.Printf("%s %s commit:%s release-build\n", appName, gitNearestTag, gitCommit)
			os.Exit(0)
		}
		fmt.Printf("%s %s commit:%s dev-build\n", appName, gitNearestTag, gitCommit)
		os.Exit(0)

	}

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

	if testcases > -1 {
		shareInfos = shareInfos[:testcases]
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

	if fillstarmodel {
		// forward data to sql
		sqldb, err := dburl.Open(fmt.Sprintf("mysql://%s:%s@%s:%d/%s", mysqlusername, mysqlpassword, mysqlhost, mysqlport, mysqldatabase))
		if err != nil {
			log.Fatal(err)
		}

		createDepartments(sqldb, flatInfos)
		createGroups(sqldb, flatInfos)
		createCompanies(sqldb, flatInfos)
		createDates(sqldb, flatInfos)
		createShares(sqldb, flatInfos)
		fmt.Println("Filled SQL star model with sharing data")
	}

	if showdata {
		if outputasjson {
			jsonData, err := json.MarshalIndent(flatInfos, "", "    ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(jsonData))
		} else {
			for _, flatInfo := range flatInfos {
				fmt.Println(flatInfo)
			}
		}
	}
	fmt.Println("Execution finished")
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
		company := models.DimensionCompany{}
		company.Company = pk
		if err := company.Insert(db); err != nil {
			log.Fatal(err)
		}
	}
}

func createGroups(db *sql.DB, flatInfos []*proto.FlatInfo) {
	groups := map[string]*proto.FlatInfo{}
	for _, flatInfo := range flatInfos {
		if _, ok := groups[flatInfo.GetOwnerInfo().Group]; !ok {
			groups[flatInfo.GetOwnerInfo().Group] = flatInfo
		}
		if _, ok := groups[flatInfo.GetShareeInfo().Group]; !ok {
			groups[flatInfo.GetShareeInfo().Group] = flatInfo
		}
	}
	for pk := range groups {
		group := models.DimensionGroup{}
		group.Egroup = pk
		if err := group.Insert(db); err != nil {
			log.Fatal(err)
		}
	}
}

func createDepartments(db *sql.DB, flatInfos []*proto.FlatInfo) {
	departments := map[string]*proto.FlatInfo{}
	for _, flatInfo := range flatInfos {
		if _, ok := departments[flatInfo.GetOwnerInfo().Department]; !ok {
			departments[flatInfo.GetOwnerInfo().Department] = flatInfo
		}
		if _, ok := departments[flatInfo.GetShareeInfo().Department]; !ok {
			departments[flatInfo.GetShareeInfo().Department] = flatInfo
		}
	}
	for pk := range departments {
		department := models.DimensionDepartment{}
		department.Department = pk
		if err := department.Insert(db); err != nil {
			log.Fatal(err)
		}
	}
}

func createDates(db *sql.DB, flatInfos []*proto.FlatInfo) {
	dates := map[int64]*proto.FlatInfo{}
	for _, flatInfo := range flatInfos {
		if _, ok := dates[flatInfo.GetShareInfo().Stime]; !ok {
			dates[flatInfo.GetShareInfo().Stime] = flatInfo
		}
	}
	for pk := range dates {
		t := time.Unix(pk, 0)
		day := t.Day()
		month := t.Month()
		year := t.Year()

		d := models.DimensionDate{}
		d.Ts = int(pk)
		d.Day = day
		d.Month = int(month)
		d.Year = year

		if err := d.Insert(db); err != nil {
			log.Fatal(err)
		}
	}
}

func createShares(db *sql.DB, flatInfos []*proto.FlatInfo) {
	for _, flatInfo := range flatInfos {
		/*
			ID               int    `json:"id"`                // id
			OwnerLogin       string `json:"owner_login"`       // owner_login
			OwnerUID         int    `json:"owner_uid"`         // owner_uid
			OwnerDepartment  string `json:"owner_department"`  // owner_department
			OwnerGroup       string `json:"owner_group"`       // owner_group
			OwnerCompany     string `json:"owner_company"`     // owner_company
			ShareeLogin      string `json:"sharee_login"`      // sharee_login
			ShareeUID        int    `json:"sharee_uid"`        // sharee_uid
			ShareeDepartment string `json:"sharee_department"` // sharee_department
			ShareeGroup      string `json:"sharee_group"`      // sharee_group
			ShareeCompany    string `json:"sharee_company"`    // sharee_company
			Stime            int    `json:"stime"`             // stime
		*/
		share := models.FactShare{}
		share.ID = int(flatInfo.GetShareInfo().Id)
		share.OwnerLogin = flatInfo.GetOwnerInfo().Login
		share.OwnerUID = int(flatInfo.GetOwnerInfo().Uid)
		share.OwnerDepartment = flatInfo.GetOwnerInfo().Department
		share.OwnerGroup = flatInfo.GetShareeInfo().Group
		share.OwnerCompany = flatInfo.GetShareeInfo().Company
		share.ShareeLogin = flatInfo.GetShareeInfo().Login
		share.ShareeUID = int(flatInfo.GetShareeInfo().Uid)
		share.ShareeDepartment = flatInfo.GetShareeInfo().Department
		share.ShareeGroup = flatInfo.GetShareeInfo().Group
		share.ShareeCompany = flatInfo.GetShareeInfo().Company
		share.Stime = int(flatInfo.GetShareInfo().Stime)

		if err := share.Insert(db); err != nil {
			log.Fatal(err)
		}
	}
}

func getPersonInfo(username string) (*proto.PersonInfo, error) {
	cmd := exec.Command(phonebookBinary, "--login", username, "-t", "login", "-t", "uid", "-t", "department", "-t", "group", "-t", "organization", "-t", "company", "-t", "office")
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
	personInfo.Group = tokens[3]
	personInfo.Organization = tokens[4]
	personInfo.Company = tokens[5]
	personInfo.Office = tokens[6]
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
