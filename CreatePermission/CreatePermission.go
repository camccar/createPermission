package CreatePermission

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"../api"
	"../models"
	_ "github.com/denisenkom/go-mssqldb"
)

func CreatePermission() {
	fmt.Println("Connection To DB and Loading Data...")

	connObj := api.GetConnection("./connection.json")

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", connObj.Server, connObj.User, connObj.Password, connObj.Database)

	//	fmt.Println(connString)

	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())

		os.Exit(1)
	} else {
		fmt.Println("Connection Sucessfully")
		fmt.Println("*********Generate Perimission Script***************")

		fmt.Println("Printing out security list")

		var list = api.GetSecurityRoles(conn)

		printList(list)
		var securityRoleId int
		fmt.Println("Which SecurtityRole should your permission fall under? Please Enter an ID from above")

		var inputstring = ""

		for inputstring == "" {
			fmt.Scan(&securityRoleId)
			if err != nil {
				fmt.Println("Problem here", nil)
			}

			inputstring = list[securityRoleId].Name
			if inputstring == "" {
				fmt.Println("Could not find an id with that input try again")
			} else {
				fmt.Println("You choose :", list[securityRoleId].Name)
			}
		}

		var list2 = api.GetSecurityActivityEnumMap(conn)

		fmt.Println("Enter a id corresponding to a section this permession should go into. a number value from below please")

		SecurityActivityMap := map[int]*models.SecurityActivity{}
		SecurityActivityMap[1] = &models.SecurityActivity{Description: "1000-1999 reserved for configuration/settings data", Label: "General"}
		SecurityActivityMap[2] = &models.SecurityActivity{Description: "2000-2999 reserved Caregiver/user related data", Label: "Careteam"}
		SecurityActivityMap[3] = &models.SecurityActivity{Description: "3000-3999 reserved for patient related data", Label: "Patient"}
		SecurityActivityMap[4] = &models.SecurityActivity{Description: "4000-4999 reserved for report related data", Label: "Reports"}
		SecurityActivityMap[5] = &models.SecurityActivity{Description: "5000-5999 reserved for API OAuth Application permissions", Label: "API OAuth"}
		SecurityActivityMap[6] = &models.SecurityActivity{Description: "6000-6999 reserved for media", Label: "Media"}
		SecurityActivityMap[7] = &models.SecurityActivity{Description: "7000-7999 reserved for authorization mode", Label: "Authorization"}

		for i := 1; i < len(SecurityActivityMap)+1; i++ {
			fmt.Println("ID: ", i, "Description:", SecurityActivityMap[i].Description)

		}

		inputstring = ""
		var inputid int
		var SecurityActivityMapId = 0
		for inputstring == "" {
			fmt.Scan(&inputid)
			if err != nil {
				fmt.Println("Problem here", nil)
			}
			inputstring = SecurityActivityMap[inputid].Description
			SecurityActivityMapId = inputid
			if inputstring == "" {
				fmt.Println("Could not find an id with that input try again")
			} else {
				fmt.Println("You chose :", SecurityActivityMap[inputid].Description)
			}
		}

		var SecurityActivityId = returnSecurityActivityNumber(inputid, list2)
		fmt.Println("Value inserting: ", SecurityActivityId)

		fmt.Println("Please choose a level of security below, low is fine in most cases")
		var SecurityLevel map[int]string
		SecurityLevel = make(map[int]string)
		SecurityLevel[1] = "Low"
		SecurityLevel[2] = "Medium - Vivify Support"
		SecurityLevel[3] = "High - Only developers"

		for i := 1; i < len(SecurityLevel)+1; i++ {
			fmt.Println("ID: ", i, " Level: ", SecurityLevel[i])
		}

		var level int

		var securityLevelChosen = 1034
		inputstring = ""
		for inputstring == "" {
			fmt.Scan(&level)
			inputstring = SecurityLevel[level]

			if inputstring == "" {
				fmt.Println("Could not find an id with that input try again")
			} else {
				fmt.Println("You chose :", SecurityLevel[level])
				if level == 1 {
					securityLevelChosen = 1034
				} else if level == 2 {
					securityLevelChosen = 1036
				} else {
					securityLevelChosen = 1038
				}
			}

		}

		fmt.Print("\nPlease enter your new permission  name with no spaces:")

		fmt.Scan(&inputstring)

		fmt.Println("\nName is: ", inputstring)

		reader := bufio.NewReader(os.Stdin)
		desctiptionnew, _ := reader.ReadString('\n')
		fmt.Print("Please Enter a description: ")
		desctiptionnew, _ = reader.ReadString('\n')

		//Not doing this anymore
		//CreateMigrateScript(list[securityRoleId].Name, SecurityActivityId, inputstring, strings.TrimSpace(desctiptionnew), securityLevelChosen, *SecurityActivityMap[SecurityActivityMapId])

		fmt.Println("\n********************Next Steps**********************\n")
		fmt.Println("Add the following line to Database/Data/dbo.SecurityActivifyEnum.Data.sql")

		fmt.Println("\n****************************************************\n")
		fmt.Printf("(%d,'%s','%s',%d,'%s')\n", SecurityActivityId, inputstring, SecurityActivityMap[SecurityActivityMapId].Label, securityLevelChosen, strings.TrimSpace(desctiptionnew))
		fmt.Println("\n****************************************************\n")
		fmt.Println("Add the following line to Database/Data/dbo.SecurityActivityRoleRoleREL.Data.sql")
		fmt.Println("\n****************************************************\n")
		fmt.Printf("(%d,%d)\n", SecurityActivityId, securityRoleId)
		fmt.Println("\n****************************************************\n")
		fmt.Println("Add the following line to vivify-platform/VivifyPlatform/Vivify.Common/Enums/SecurityActivityEnum.cs")
		fmt.Println("\n****************************************************\n")
		fmt.Printf("%s = %d\n", inputstring, SecurityActivityId)
		fmt.Println("\n****************************************************")
		fmt.Println("Add the following line to vivify-platform/Database/Migrations/Config/Permissions.sql")
		fmt.Println("\n****************************************************\n")
		fmt.Printf("(%d,@RoleId)\n\n", SecurityActivityId)
		fmt.Println("\n****************************************************\n")
		fmt.Println("Add the following line to Vivify.Platform/Database/Migrations/Config/ActivityEnum.sql")
		fmt.Println("\n****************************************************\n")
		fmt.Printf("(%d,'%s','%s',%d,'%s')\n\n", SecurityActivityId, inputstring, strings.TrimSpace(desctiptionnew), securityLevelChosen, SecurityActivityMap[SecurityActivityMapId].Label)
		fmt.Println("\n****************************************************\n")
		fmt.Println("Add the following line to the front end code in vh-web-caregiverportal/src/modules/permissions/permissions.service.js")
		fmt.Println("\n****************************************************\n")
		fmt.Printf("%s:%d,\n\n", inputstring, SecurityActivityId)
		fmt.Println("\n****************************************************\n")
		var exitinput = ""
		fmt.Printf("Finished Script type anything and press enter to end application\n")
		fmt.Scan(&exitinput)
		fmt.Printf("Ending Application\n")

	}
	defer conn.Close()
}

func CreateMigrateScript(SecurityRole string, id int, name string, description string, securityLevelChosen int, SecurityActivityMap models.SecurityActivity) {

	output := fmt.Sprintf("IF NOT EXISTS (SELECT 1 FROM SecurityActivityEnum Where SecurityActivityId =  %d )\n    Begin\n", id)
	output = output + fmt.Sprintf("        INSERT INTO SecurityActivityEnum(SecurityActivityId, Name, Description, FilterSecurityActivityId,\"Group\")\n")
	output = output + fmt.Sprintf("        VALUES ( %d ,'%s' ,'%s',%d,'%s'  )\n", id, name, description, securityLevelChosen, SecurityActivityMap.Label)
	output = output + fmt.Sprintf("    End\n")
	output = output + fmt.Sprintf("IF NOT EXISTS (SELECT 1 FROM SecurityActivityRoleRel WHERE SecurityActivityId = %d AND SecurityRoleId = (SELECT SecurityRoleId FROM SecurityRole WHERE Name = '%s'))\n", id, SecurityRole)
	output = output + fmt.Sprintf("BEGIN\n")
	output = output + fmt.Sprintf("    INSERT INTO SecurityActivityRoleREL VALUES (%d, ( SELECT SecurityRoleId FROM SecurityRole WHERE Name = '%s'))\n", id, SecurityRole)
	output = output + fmt.Sprintf("END\n")

	filename := fmt.Sprintf("IHCP-____Add%sPermission.sql", name)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Something bad happened", err)
	} else {
		fmt.Println("Writting migration script to ", filename)
	}
	defer file.Close()

	fmt.Fprintf(file, output)

}

func returnSecurityActivityNumber(section int, input map[int]models.SecurityActivityEnum) int {

	var counter = 0

	switch section {
	case 1:
		counter = 1000
	case 2:
		counter = 2000
	case 3:
		counter = 3000
	case 4:
		counter = 4000
	case 5:
		counter = 5000
	case 6:
		counter = 6000
	case 7:
		counter = 7000
	default:
		counter = 1000

	}
	var returnval = 0

	for begin := counter; begin < counter+1000; begin++ {
		if input[begin].Name == "" {
			returnval = begin
			break
		}
	}

	return returnval

}

func printList(inputL map[int]models.Security) {

	for i := 1; i < len(inputL)+1; i++ {
		if inputL[i].Id > 9 {
			if inputL[i].Id != 0 {
				fmt.Println("Id:", inputL[i].Id, " Name:", inputL[i].Name, "                  ")
			}
		} else {
			if inputL[i].Id != 0 {
				fmt.Println("Id:", inputL[i].Id, "  Name:", inputL[i].Name, "                  ")
			}
		}
	}

}
