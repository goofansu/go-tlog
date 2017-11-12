package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// XMLEntry 对应XML中的entry
type XMLEntry struct {
	XMLName xml.Name `xml:"entry"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Desc    string   `xml:"desc,attr"`
}

// XMLStruct 对应XML中的struct
type XMLStruct struct {
	XMLName xml.Name   `xml:"struct"`
	Name    string     `xml:"name,attr"`
	Desc    string     `xml:"desc,attr"`
	Entries []XMLEntry `xml:"entry"`
}

// XMLMetalib 对应XML中的metalib
type XMLMetalib struct {
	XMLName xml.Name    `xml:"metalib"`
	Structs []XMLStruct `xml:"struct"`
}

var source string
var dest string
var pkg string

func init() {
	flag.StringVar(&source, "source", "tlog.xml", "XML file path, support glob pattern like */**.xml")
	flag.StringVar(&dest, "dest", "tlogevt.go", "Go file built from XMLs")
	flag.StringVar(&pkg, "pkg", "tlogevt", "Go file package name")
}

func main() {
	flag.Parse()

	f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	xmlFilePaths, _ := filepath.Glob(source)
	xmlStructs := allStructs(xmlFilePaths)
	sort.Slice(xmlStructs, func(i, j int) bool {
		return xmlStructs[i].Name < xmlStructs[j].Name
	})

	fmt.Fprintf(f, "// Generated at: %s\n", time.Now())
	fmt.Fprint(f, "// Auto-generated from XMLs, don't modify.\n")
	fmt.Fprintf(f, "package %s\n", pkg)

	for _, s := range xmlStructs {
		writeStruct(f, s)
	}
}

func allStructs(xmlFilePaths []string) []XMLStruct {
	var result []XMLStruct
	for _, fpath := range xmlFilePaths {
		file, err := os.Open(fpath)
		if err != nil {
			log.Fatalln(err, fpath)
		}
		defer file.Close()

		ss, err := readStructs(file)
		if err != nil {
			log.Fatalln(err, file)
		}
		result = append(result, ss...)
	}
	return result
}

func writeStruct(w io.Writer, s XMLStruct) {
	fmt.Fprintf(w, "// %s %s\n", s.Name, s.Desc)
	fmt.Fprintf(w, "type %s struct {\n", s.Name)
	for _, e := range s.Entries {
		required := isFieldRequired(e.Desc)
		fmt.Fprintf(w, "%s %s %s // %s\n",
			fieldName(e.Name),
			fieldType(e.Type, e.Name, required),
			fieldTag(e.Name, required),
			e.Desc)
	}
	fmt.Fprint(w, "}\n")
}

func readStructs(reader io.Reader) ([]XMLStruct, error) {
	var xmlMetalib XMLMetalib
	if err := xml.NewDecoder(reader).Decode(&xmlMetalib); err != nil {
		return nil, err
	}
	return xmlMetalib.Structs, nil
}

func fieldName(name string) string {
	re := regexp.MustCompile("(Id$)|(^cpu)")
	s := strings.Title(name)
	return re.ReplaceAllStringFunc(s, func(w string) string {
		return strings.ToUpper(w)
	})
}

func fieldType(ttype string, name string, required bool) string {
	ttype = _fieldType(ttype, name)
	if required {
		return "*" + ttype
	}
	return ttype
}

func _fieldType(ttype string, name string) string {
	switch ttype {
	case "datetime":
		return "time.Time"
	case "float":
		return "float64"
	case "int":
		if name == "Sequence" {
			return "int64"
		}
		return ttype
	default:
		return ttype
	}
}

func fieldTag(name string, required bool) string {
	if required {
		return fmt.Sprintf("`tlog:\"%s\" validate:\"required\"`", name)
	}
	return fmt.Sprintf("`tlog:\"%s\"`", name)
}

func isFieldRequired(desc string) bool {
	return strings.Index(desc, "必填") != -1
}
