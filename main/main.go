package main

import "../service"

const USER_ID = "1718998874"
const Dirname = `D:\\weiboBackup\`

func main() {
    service.GetAllWeibo(USER_ID, Dirname)
}