package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

func IsEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false // Either not empty or error, suits both cases
}

func rdirTree(dir string, lvl int, folders bool, out *bytes.Buffer, arr [50]int) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	arr[lvl] = 1

	cntfld := 0
	endlvl := 0

	enter := false

	for i := 0; i <= lvl; i++ {
		if arr[i] == 0 {
			enter = true
		}
	}

	for i, file := range files {
		if file.IsDir() == true {
			cntfld = i
		}
	}

	for i, file := range files {
		if file.IsDir() == false {
			if folders {
				if lvl == 1 {
					out.WriteString("	")
				} else {
					if lvl != 0 {
						out.WriteString("	")
						spaces := (lvl - 1) * 2
						for cnt := 1; cnt <= spaces; cnt++ {
							if arr[(cnt/2)+1] == 1 {
								if (cnt % 2) == 1 {
									out.WriteString("│")
								} else {
									out.WriteString("	")
								}
							} else {
								if (cnt % 2) == 1 {
									out.WriteString("")
								} else {
									out.WriteString("	")
								}
							}
						}

					}
				}
				if file.Name() != "main.go" {
					if i == len(files)-1 {
						endlvl++
						if int(file.Size()) > 0 {
							arr[lvl] = 0
							if enter == false {
								out.WriteString("└───" + file.Name() + " (" + strconv.Itoa(int(file.Size())) + "b)")
							} else {
								out.WriteString("└───" + file.Name() + " (" + strconv.Itoa(int(file.Size())) + "b)" + "\n")
							}
						} else {
							arr[lvl] = 0
							if enter == false {
								out.WriteString("└───" + file.Name() + " (empty)")
							} else {
								out.WriteString("└───" + file.Name() + " (empty)" + "\n")
							}
						}
					} else {
						if int(file.Size()) > 0 {
							out.WriteString("├───" + file.Name() + " (" + strconv.Itoa(int(file.Size())) + "b)" + "\n")
						} else {
							out.WriteString("├───" + file.Name() + " (empty)" + "\n")
						}
					}
				} else {
					if i == len(files)-1 {
						arr[lvl] = 0

						if enter == false {
							out.WriteString("└───" + file.Name() + " (" + "vary" + ")")
						} else {
							out.WriteString("└───" + file.Name() + " (" + "vary" + ")" + "\n")
						}

					} else {
						out.WriteString("├───" + file.Name() + " (" + "vary" + ")" + "\n")
					}
				}
			}
		} else {
			if lvl == 1 {
				out.WriteString("	")
			} else {
				if lvl != 0 {
					out.WriteString("	")
					spaces := (lvl - 1) * 2
					for cnt := 1; cnt <= spaces; cnt++ {
						if arr[(cnt/2)+1] == 1 {
							if (cnt % 2) == 1 {
								out.WriteString("│")
							} else {
								out.WriteString("	")
							}
						} else {
							if (cnt % 2) == 1 {
								out.WriteString("")
							} else {
								out.WriteString("	")
							}
						}
					}

				}
			}

			if true {

				if (i == len(files)-1 && folders == true) || (i == cntfld && folders == false) {
					arr[lvl] = 0

					if enter == false && IsEmpty(dir+"/"+file.Name()) {
						out.WriteString("└───" + file.Name())
					} else {
						out.WriteString("└───" + file.Name() + "\n")
					}
				} else {
					out.WriteString("├───" + file.Name() + "\n")
				}

				newlvl := lvl + 1
				newdir := dir + "/" + file.Name()
				rdirTree(newdir, newlvl, folders, out, arr)
			}
		}
	}

	return nil
}

func dirTree(out *bytes.Buffer, path string, printFiles bool) error {
	var array [50]int
	for i := 0; i < 50; i++ {
		array[i] = 0
	}
	return rdirTree(path, 0, printFiles, out, array)
}

func main() {
	out := new(bytes.Buffer)
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	print(out.String())
	if err != nil {
		panic(err.Error())
	}
}
