package cmd

import (
	"bufio"
	"fmt"
	"github.com/newlee/tequila/viz"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var javaDbCmd *cobra.Command = &cobra.Command{
	Use:   "jd",
	Short: "java code to database dependencies",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		source := cmd.Flag("source").Value.String()
		filter := cmd.Flag("filter").Value.String()
		pkgKey := cmd.Flag("package").Value.String()
		tableKey := cmd.Flag("table").Value.String()

		codeFiles := make([]string, 0)
		filepath.Walk(source, func(path string, fi os.FileInfo, err error) error {
			if !strings.HasSuffix(path, ".java") {
				return nil
			}

			split := strings.Split(filter, ",")
			if filter == "" {
				codeFiles = append(codeFiles, path)
			} else {
				for _, key := range split {
					if strings.Contains(path, key) {
						return nil
					}
				}
				codeFiles = append(codeFiles, path)
			}

			return nil
		})

		allPkg := viz.NewAllPkg()
		tables := viz.NewAllTable()
		pkgCallerFiles := make(map[string]string)
		tableCallerFiles := make(map[string]string)

		var pkgFilter func(line string) bool
		if filter == "" {
			pkgFilter = func(line string) bool {
				return !strings.HasPrefix(line, pkgKey)
			}
		} else {
			pkgFilter = func(line string) bool {
				return strings.HasPrefix(line, pkgKey)
			}
		}

		var tableFilter func(line string) bool
		if filter == "" {
			tableFilter = func(line string) bool {
				return !strings.HasPrefix(line, tableKey)
			}
		} else {
			tableFilter = func(line string) bool {
				return strings.HasPrefix(line, tableKey)
			}
		}
		for _, codeFileName := range codeFiles {
			codeFile, _ := os.Open(codeFileName)
			scanner := bufio.NewScanner(codeFile)
			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				line := scanner.Text()
				line = strings.ToUpper(line)
				fields := strings.Fields(line)
				if len(fields) == 0 {
					continue
				}
				first := fields[0]
				if strings.HasPrefix(first, "/*") || strings.HasPrefix(first, "*") || strings.HasPrefix(first, "//") {
					continue
				}
				if strings.Contains(line, "PKG_") {
					tmp := strings.FieldsFunc(line, func(r rune) bool {
						return r == ' ' || r == '(' || r == ',' || r == '\''
					})

					for _, key := range tmp {
						if strings.HasPrefix(key, "PKG_") && pkgFilter(key) {
							sk := strings.Replace(key, "\"", "", -1)
							spss := strings.Split(sk, ".")
							pkg := spss[0]

							if len(spss) > 1 && strings.Contains(pkg, "PKG_") {
								split := strings.Split(codeFileName, "/")
								pkgCallerFiles[split[len(split)-1]] = ""
								sp := spss[1]
								if strings.HasPrefix(sp, "P_") {
									allPkg.Add(pkg, sp)
								}
							}
						}
					}
				}

				if strings.Contains(line, " T_") || strings.Contains(line, ",T_") {
					tmp := strings.FieldsFunc(line, func(r rune) bool {
						return r == ' ' || r == ',' || r == '.' || r == '"' || r == ':' || r == '(' || r == ')' || r == '）' || r == '%' || r == '!' || r == '\''
					})
					for _, t3 := range tmp {
						if strings.HasPrefix(t3, "T_") && tableFilter(t3) && !viz.IsChineseChar(t3) {
							split := strings.Split(codeFileName, "/")
							tableCallerFiles[split[len(split)-1]] = ""
							tables.Add(t3)
						}
					}
				}
			}

			codeFile.Close()
		}

		allPkg.Print()
		fmt.Println("")
		fmt.Println("-----------")
		for key := range pkgCallerFiles {
			fmt.Println(key)

		}
		fmt.Println("")
		fmt.Println("-----------")

		tables.Print()

		fmt.Println("")
		fmt.Println("-----------")
		for key := range tableCallerFiles {
			fmt.Println(key)

		}
	},
}

func init() {
	rootCmd.AddCommand(javaDbCmd)

	javaDbCmd.Flags().StringP("source", "s", "", "source code directory")
	javaDbCmd.Flags().StringP("filter", "f", "", "file filter")
	javaDbCmd.Flags().StringP("table", "t", "T_", "table filter")
	javaDbCmd.Flags().StringP("package", "p", "PKG_", "package filter")
	javaDbCmd.Flags().StringP("output", "o", "java_db.dot", "output dot file name")
}
