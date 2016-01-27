package main
import (
    "fmt"
	"log"
    "os"
    "path/filepath"
    "io/ioutil"
    "strings"
    "encoding/base64"
)



func main(){
	if len(os.Args) != 4 {
		log.Fatal("Impossible to autogenerate resources, miss arguments. Usage : [resources folder] [package] [target resources folder]")
	}
    inputResourcesFolder := os.Args[1]
    // Insert generator files with data into specific package
    packageFolder := os.Args[2]
    resourcesFolder := os.Args[3]
	log.Println("Create autogenerate for ",inputResourcesFolder,"in package",packageFolder,"with target",resourcesFolder)
    outPath := filepath.Join(packageFolder,"autogenerate_resources.go")
	outFile,_ := os.OpenFile(outPath,os.O_CREATE|os.O_RDWR|os.O_TRUNC,os.ModePerm)
    outFile.WriteString("package " + packageFolder + "\n")
    outFile.WriteString("import \"os\"\n")
    outFile.WriteString("import \"strings\"\n")
    outFile.WriteString("import \"path/filepath\"\n")
    outFile.WriteString("import \"io/ioutil\"\n")
    outFile.WriteString("import \"encoding/base64\"\n")
    outFile.WriteString("import \"fmt\"\n\n")

    outFile.WriteString("var resourcesFolder = \"" + resourcesFolder + "\"\n\n")

    outFile.WriteString("func init(){\n")
    outFile.WriteString("files := map[string]string{\"\":\"\"")

    treat(outFile,inputResourcesFolder,"")
    outFile.WriteString("}\n\n")

    writeCode(outFile,resourcesFolder)
    outFile.WriteString("}\n")
    outFile.Close()
	log.Println("Code generate in file",outPath)
}

func writeCode(out *os.File,resourcesFolder string){
    out.WriteString("for _,a:= range os.Args[1:]{\n")
    out.WriteString("\tif a == \"-forceDeploy\"{\n")
    out.WriteString("\t\tfmt.Println(\"Force remove folder\")\n")
    out.WriteString("\t\tos.RemoveAll(resourcesFolder)\n")
    out.WriteString("\t\tbreak\n")
    out.WriteString("\t}\n")
    out.WriteString("}\n\n")

    out.WriteString("if _,err := os.Open(resourcesFolder) ; err == nil{\n")
    out.WriteString("\treturn\n")
    out.WriteString("}\n")
    out.WriteString("\n\nfmt.Println(\"Generating resources in folder\",resourcesFolder)\n")
    out.WriteString("for name,data := range files {\n")
    out.WriteString("\tif name!=\"\" {\n")
    out.WriteString("\t\td:=resourcesFolder\n")
    sep := ""
    if os.PathSeparator == '\\' {
        sep = fmt.Sprintf("\\%c",os.PathSeparator)
    }else {
        sep = fmt.Sprintf("%c",os.PathSeparator)
    }

    out.WriteString(fmt.Sprintf("\t\tif idx:= strings.LastIndex(name,\"%s\") ; idx !=-1 {\n",sep))
    out.WriteString("\t\t\td=filepath.Join(d,name[:idx])\n")
    out.WriteString("\t\t}\n")
    out.WriteString("\t\tos.MkdirAll(d,os.ModePerm)\n")
    out.WriteString("\t\tdecodeData,_ := base64.StdEncoding.DecodeString(data)\n")
    out.WriteString("\t\tioutil.WriteFile(filepath.Join(resourcesFolder,name),decodeData,os.ModePerm)\n")
    out.WriteString("\t\tfmt.Println(\"=>\",d,\":\",name,len(decodeData))\n")
    out.WriteString("\t}\n}\n")

    out.WriteString("\n")
}

func treat(outFile *os.File,root,dir string){
    f,_ := os.Open(filepath.Join(root,dir))
    files,_ := f.Readdir(-1)

    //r2,_ := regexp.Compile("//.*\r\n")
    for _,file := range files {
        if file.IsDir() {
            dirName := filepath.Join(dir,file.Name())
            treat(outFile,root,dirName)
        }else{
            in := filepath.Join(root,dir,file.Name())
            data,_ := ioutil.ReadFile(in)
            log.Println("Add file",in)
            strData := base64.StdEncoding.EncodeToString(data)
            outFile.WriteString(",\"" + strings.Replace(filepath.Join(dir,file.Name()),"\\","\\\\",-1) + "\":`" + strData  + "`")
        }
    }
}
