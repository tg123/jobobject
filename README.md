# JobObject
 Windows JobObject utils for kill all child processes when parent process exits
 
# Exmaple

```
func main() {
	e := exec.Command("notepad.exe")
	e.Start()

	job, err := jobobject.Create()
	if err != nil {
		panic(err)
	}

	err = job.AddProcess(e.Process)
	if err != nil {
		panic(err)
	}

	// defer job.Close() // This jobobj will be closed even if the parent process is killed
	fmt.Println("ctrl + c to quit and kill notepad.exe")
	time.Sleep(time.Hour * 1)
 }
```
