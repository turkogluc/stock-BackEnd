package controllers

import (
	"github.com/gin-gonic/gin"
	."stock/entities"
	"strconv"
	"strings"
	"stock/entities/responses"
	. "stock/common/logger"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"os"
	"io"
)

var imagePath string

func SetImagePath(path string){
	imagePath = path
}

func InitRoutes(public,private *gin.RouterGroup) {

	public.POST("login", handleLogin)
	private.GET("me", handleGetMe)
	private.GET("fillProductTable", fillProductTable)

	private.POST("createProduct", createProduct)
	private.POST("updateProduct", updateProduct)
	private.GET("getProductById", getProductById)
	private.GET("deleteProducts", deleteProducts)
	private.GET("getProducts",getProducts)
	private.GET("retrieveCategories",retrieveCategories)
	private.POST("uploadFile", uploadFile)

	private.POST("createStock", createStock)
	private.POST("updateStock", updateStock)
	private.GET("getStockById", getStockById)
	private.GET("deleteStocks", deleteStocks)
	private.GET("getStocks",getStocks)
	private.GET("setFavoriteProduct", setFavoriteProduct)

	private.POST("createPerson", createPerson)
	private.POST("updatePerson", updatePerson)
	private.GET("getPersonById", getPersonById)
	private.GET("deletePeople", deletePeople)
	private.GET("getPeople",getPeople)

	private.POST("createReceiving", createReceiving)
	private.POST("updateReceiving", updateReceiving)
	private.GET("getReceivingById", getReceivingById)
	private.GET("deleteReceivings", deleteReceivings)
	private.GET("getReceivings",getReceivings)
	private.GET("setReceivingStatus", setReceivingStatus)

	private.POST("createPayment", createPayment)
	private.POST("updatePayment", updatePayment)
	private.GET("getPaymentById", getPaymentById)
	private.GET("deletePayments", deletePayments)
	private.GET("getPayments",getPayments)
	private.GET("setPaymentStatus", setPaymentStatus)

	private.POST("createExpense", createExpense)
	private.POST("updateExpense", updateExpense)
	private.GET("getExpenseById", getExpenseById)
	private.GET("deleteExpenses", deleteExpenses)
	private.GET("getExpenses",getExpenses)

	public.POST("createUser", createUser)
	private.POST("updateUser", updateUser)
	private.GET("getUserById", getUserById)
	private.GET("deleteUsers", deleteUsers)
	private.GET("getUsers",getUsers)

	private.POST("createSale", createSale)
	private.POST("updateSale", updateSale)
	private.GET("getSaleById", getSaleById)
	private.GET("deleteSales", deleteSales)
	private.GET("getSales",getSales)

	// # Reports #

	private.GET("getSaleSummaryReport", getSaleSummaryReport)
	private.GET("getCurrentStockReport",getCurrentStockReport)
	private.GET("retrieveActivityLog",getActivityLog)
	private.GET("getPaymentReport",getPaymentReport)
	private.GET("getProductReport",getProductReport)

	// # Excel Reports #

	private.GET("getSaleSummaryReportAsExcel", getSaleSummaryReportAsExcel)
	private.GET("getCurrentStockReportAsExcel",getCurrentStockReportAsExcel)
	private.GET("getPaymentReportAsExcel",getPaymentReportAsExcel)
	private.GET("getProductReportAsExcel",getProductReportAsExcel)

}

func handleGetMe(c *gin.Context){
	var err *ErrorType

	uId := getUserIdFromToken(c)

	ur, err := UseCase.GetUserById(uId)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(LoginResponse{User: ur}))

}

func fillProductTable(c *gin.Context){

	userId := getUserIdFromToken(c)
	err := UseCase.FillProductTable(userId)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))

}

func getUserIdFromToken(c *gin.Context) int {
	v, _ := c.Get("token-claims")
	LogDebug(v)
	claims := v.(jwt.MapClaims)
	userId := claims["userId"].(float64)
	return int(userId)
}


func createProduct (c *gin.Context){
	p := Product{}
	c.BindJSON(&p)

	//// upload picture
	//file, err := c.FormFile("file")
	//if err != nil && file != nil {
	//	c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
	//	return
	//}else{
	//	fileName := strconv.Itoa(p.Id) + "jpg"
	//
	//	if err := c.SaveUploadedFile(file, fileName); err != nil {
	//		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	//		return
	//	}
	//}

	p.UserId = getUserIdFromToken(c)
	err2 := UseCase.CreateProduct(&p)

	if err2 != nil{
		c.JSON(200, generateFailResponse(err2))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func updateProduct (c *gin.Context){
	p := Product{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.UpdateProduct(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getProductById (c *gin.Context){
	id,_ := strconv.Atoi(c.Query("id"))

	p, err := UseCase.GetProductById(id)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func deleteProducts(c *gin.Context){

	idList := strings.Split(c.Query("ids"),",")

	var ids []int
	for _,id := range idList{
		i,_ := strconv.Atoi(id)
		ids = append(ids,i)
	}

	err := UseCase.DeleteProducts(ids)
	if err != nil {
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

func getProducts (c *gin.Context){

	barcode := c.Query("barcode")
	name := c.Query("name")
	description := c.Query("description")
	category := c.Query("category")

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")
	isDropdown,_ := strconv.ParseBool(c.Query("isDropdown"))

	p, err := UseCase.GetProducts(barcode,name,description,category,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	if isDropdown {
		ResponseList := []*responses.ProductDropdownResponse{}
		for _,product := range p.Items{
			r := &responses.ProductDropdownResponse{}
			r.Id = product.Id
			r.Name = product.Name
			r.Price = product.SalePrice
			ResponseList = append(ResponseList, r)
		}
		c.JSON(200, generateSuccessResponse(ResponseList))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func retrieveCategories (c *gin.Context){

	res,err := UseCase.RetrieveCategories()
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(res))

}

func uploadFile(c *gin.Context) {

	//pId := getProfileIdFromToken(c)

	file, header , err := c.Request.FormFile("file")
	if err != nil{
		fmt.Println(err)
		c.JSON(200, nil)
	}
	filename := header.Filename


	//filename = strconv.Itoa(int(time.Now().Unix()))

	if err != nil{
		LogError(err)
		c.JSON(200, err)
		return
	}

	// image path is /var/www/html/ximage/
	out, err := os.Create(imagePath + filename)
	if err != nil {
		LogError(err)
		c.JSON(200, err)
		return
	}

	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		LogError(err)
		c.JSON(200, nil)
	}

	ipAddress := "http://128.199.53.5/ximages/" + filename

	c.JSON(200,gin.H{"url": ipAddress})
}

func removeFile(c *gin.Context) {

	path := c.Query("path")

	err := os.Remove(path)
	if err != nil {
		LogError(err)
		c.JSON(200, err)
		return
	}
	c.JSON(200, nil)

}

//func DownloadMedia(c *gin.Context) {
//	m := MediaURL{}
//	err := c.BindJSON(&m)
//	if err != nil {
//		LogError(err)
//		c.JSON(200, nil)
//		return
//	}
//
//
//
//	//c.Header("Content-Encoding","utf-8")
//	//c.Header("Content-Description", "File Transfer")
//	//c.Header("Content-Transfer-Encoding", "binary")
//	c.Header("Content-Disposition", `attachment;filename="` + m.FileName + `"`)
//	c.Header("Content-Type",  m.Type)
//	c.Header("Content-Length", strconv.FormatInt(int64(m.Size), 10))
//	c.File("/home/iugo/WebFleet" + m.Url)
//
//
//	//c.Header("Content-Description", "File Transfer")
//	//c.Header("Content-Disposition", `attachment; filename=TripMeasData.text` )
//	//c.Header("Content-Type", "application/octet-stream")
//	//c.File("/home/iugo/WebFleet/media/TripMeasData_app.txt")
//
//}

// ########################################################

func createStock (c *gin.Context){
	p := Stock{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.CreateStock(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func updateStock (c *gin.Context){
	p := Stock{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.UpdateStock(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getStockById (c *gin.Context){
	id,_ := strconv.Atoi(c.Query("id"))

	p, err := UseCase.GetStockById(id)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func deleteStocks(c *gin.Context){

	idList := strings.Split(c.Query("ids"),",")

	var ids []int
	for _,id := range idList{
		i,_ := strconv.Atoi(id)
		ids = append(ids,i)
	}

	err := UseCase.DeleteStocks(ids)
	if err != nil {
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

func getStocks (c *gin.Context){

	tInterval := c.Query("timeInterval")

	barcode := c.Query("barcode")
	name := c.Query("name")
	description := c.Query("description")
	category := c.Query("category")
	dealerId,_ := strconv.Atoi(c.Query("dealerId"))
	creator,_ := strconv.Atoi(c.Query("creatorId"))

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")
	//isDropdown,_ := strconv.ParseBool(c.Query("isDropdown"))
	isFavorite,_ := strconv.ParseBool(c.Query("isFavorite"))

	p, err := UseCase.GetStocks(tInterval,barcode,name,description,category,orderBy,orderAs,pageNumber, pageSize,dealerId,creator,isFavorite)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func setFavoriteProduct (c *gin.Context){

	id,_ := strconv.Atoi(c.Query("productId"))
	isFavorite,_ := strconv.ParseBool(c.Query("isFavorite"))

	err := UseCase.SetFavoriteProduct(id,isFavorite)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

// ###########################################

func createPerson (c *gin.Context){
	p := Person{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.CreatePerson(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func updatePerson (c *gin.Context){
	p := Person{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.UpdatePerson(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getPersonById (c *gin.Context){
	id,_ := strconv.Atoi(c.Query("id"))

	p, err := UseCase.GetPersonById(id)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func deletePeople(c *gin.Context){

	idList := strings.Split(c.Query("ids"),",")

	var ids []int
	for _,id := range idList{
		i,_ := strconv.Atoi(id)
		ids = append(ids,i)
	}

	err := UseCase.DeletePersons(ids)
	if err != nil {
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

func getPeople (c *gin.Context){

	name := c.Query("name")
	pType := c.Query("pType")

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")
	isDropdown,_ := strconv.ParseBool(c.Query("isDropdown"))

	p, err := UseCase.GetPeople(name,pType,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	if isDropdown {
		ResponseList := []*responses.PersonDropdownResponse{}
		for _,per := range p.Items{
			r := &responses.PersonDropdownResponse{}
			r.Id = per.Id
			r.Name = per.Name
			r.Type = per.Type
			ResponseList = append(ResponseList, r)
		}
		c.JSON(200, generateSuccessResponse(ResponseList))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

// #############################################################

func createReceiving (c *gin.Context){
	p := Receiving{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.CreateReceiving(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func updateReceiving (c *gin.Context){
	p := Receiving{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.UpdateReceiving(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getReceivingById (c *gin.Context){
	id,_ := strconv.Atoi(c.Query("id"))

	p, err := UseCase.GetReceivingById(id)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func deleteReceivings(c *gin.Context){

	idList := strings.Split(c.Query("ids"),",")

	var ids []int
	for _,id := range idList{
		i,_ := strconv.Atoi(id)
		ids = append(ids,i)
	}

	err := UseCase.DeleteReceivings(ids)
	if err != nil {
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

func getReceivings (c *gin.Context){

	person := c.Query("person")
	status := c.Query("status")

	tInterval := c.Query("timeInterval")
	creator,_ := strconv.Atoi(c.Query("creatorId"))

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")

	p, err := UseCase.GetReceivings(tInterval,person,status,orderBy,orderAs,pageNumber, pageSize,creator)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func setReceivingStatus (c *gin.Context){

	id,_ := strconv.Atoi(c.Query("id"))
	status := c.Query("status")

	err := UseCase.SetReceivingStatus(status,id)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

// #######################################################

func createPayment (c *gin.Context){
	p := Payment{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.CreatePayment(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func updatePayment (c *gin.Context){
	p := Payment{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.UpdatePayment(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getPaymentById (c *gin.Context){
	id,_ := strconv.Atoi(c.Query("id"))

	p, err := UseCase.GetPaymentById(id)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func deletePayments(c *gin.Context){

	idList := strings.Split(c.Query("ids"),",")

	var ids []int
	for _,id := range idList{
		i,_ := strconv.Atoi(id)
		ids = append(ids,i)
	}

	err := UseCase.DeletePayments(ids)
	if err != nil {
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

func getPayments (c *gin.Context){

	person := c.Query("person")
	status := c.Query("status")

	tInterval := c.Query("timeInterval")
	creator,_ := strconv.Atoi(c.Query("creatorId"))

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")

	p, err := UseCase.GetPayments(tInterval,person,status,orderBy,orderAs,pageNumber, pageSize,creator)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func setPaymentStatus (c *gin.Context){

	id,_ := strconv.Atoi(c.Query("id"))
	status := c.Query("status")

	err := UseCase.SetPaymentStatus(status,id)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

// ###############################################################


func createExpense (c *gin.Context){
	p := Expense{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.CreateExpense(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func updateExpense (c *gin.Context){
	p := Expense{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.UpdateExpense(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getExpenseById (c *gin.Context){
	id,_ := strconv.Atoi(c.Query("id"))

	p, err := UseCase.GetExpenseById(id)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func deleteExpenses(c *gin.Context){

	idList := strings.Split(c.Query("ids"),",")

	var ids []int
	for _,id := range idList{
		i,_ := strconv.Atoi(id)
		ids = append(ids,i)
	}

	err := UseCase.DeleteExpenses(ids)
	if err != nil {
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

func getExpenses (c *gin.Context){

	name := c.Query("name")
	description := c.Query("description")

	tInterval := c.Query("timeInterval")
	creator,_ := strconv.Atoi(c.Query("creatorId"))

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")
	//isDropdown,_ := strconv.ParseBool(c.Query("isDropdown"))

	p, err := UseCase.GetExpenses(tInterval,name,description,orderBy,orderAs,pageNumber, pageSize,creator)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

// ###############################################################

func createUser (c *gin.Context){
	p := User{}
	c.BindJSON(&p)

	err := UseCase.CreateUser(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func updateUser (c *gin.Context){
	p := User{}
	c.BindJSON(&p)

	err := UseCase.UpdateUser(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getUserById (c *gin.Context){
	id,_ := strconv.Atoi(c.Query("id"))

	p, err := UseCase.GetUserById(id)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func deleteUsers(c *gin.Context){

	idList := strings.Split(c.Query("ids"),",")

	var ids []int
	for _,id := range idList{
		i,_ := strconv.Atoi(id)
		ids = append(ids,i)
	}

	err := UseCase.DeleteUsers(ids)
	if err != nil {
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

func getUsers (c *gin.Context){

	name := c.Query("name")
	email := c.Query("email")

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")
	isDropdown,_ := strconv.ParseBool(c.Query("isDropdown"))

	p, err := UseCase.GetUsers(name,email,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	if isDropdown {
		ResponseList := []*responses.UserDropdownResponse{}
		for _,per := range p.Items{
			r := &responses.UserDropdownResponse{}
			r.Id = per.Id
			r.Name = per.Name
			r.Email = per.Email
			r.Phone = per.Phone
			ResponseList = append(ResponseList, r)
		}
		c.JSON(200, generateSuccessResponse(ResponseList))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

// #########################################################

func createSale (c *gin.Context){
	p := SaleBasket{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.CreateSaleBasket(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func updateSale (c *gin.Context){
	p := SaleBasket{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.UpdateSaleBasket(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getSaleById (c *gin.Context){
	id,_ := strconv.Atoi(c.Query("id"))

	p, err := UseCase.GetSaleBasketById(id)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func deleteSales(c *gin.Context){

	idList := strings.Split(c.Query("ids"),",")

	var ids []int
	for _,id := range idList{
		i,_ := strconv.Atoi(id)
		ids = append(ids,i)
	}

	err := UseCase.DeleteSaleBaskets(ids)
	if err != nil {
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

func getSales (c *gin.Context){

	tInterval := c.Query("timeInterval")
	userId,_ := strconv.Atoi(c.Query("userId"))

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")
	//isDropdown,_ := strconv.ParseBool(c.Query("isDropdown"))

	p, err := UseCase.GetSaleBaskets(tInterval,userId,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

// ###############################################################

// # reports #

func getSaleSummaryReport(c *gin.Context){

	tInterval := c.Query("timeInterval")

	p, err := UseCase.GetSaleSummaryReport(tInterval)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getCurrentStockReport(c *gin.Context){

	//barcode := c.Query("barcode")
	name := c.Query("name")
	category := c.Query("category")

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")

	p, err := UseCase.GetCurrentStockReport(name,category,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{	
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getActivityLog (c *gin.Context){

	tInterval := c.Query("timeInterval")
	userId,_ := strconv.Atoi(c.Query("userId"))

	p, err := UseCase.GetActivityLog(tInterval,userId)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getPaymentReport (c *gin.Context){

	tInterval := c.Query("timeInterval")

	p, err := UseCase.GetPaymentReport(tInterval)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getProductReport (c *gin.Context){

	tInterval := c.Query("timeInterval")
	category := c.Query("category")
	productName := c.Query("productName")
	userId,_ := strconv.Atoi(c.Query("userId"))


	p, err := UseCase.GetProductReport(tInterval,productName,category,userId)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}


// # reports to excel #

func getSaleSummaryReportAsExcel(c *gin.Context){


	tInterval := c.Query("timeInterval")

	fileName, err := UseCase.GetSaleSummaryReportAsExcel(tInterval)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	println("file:",fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", `attachment; filename=excelFile.xlsx` )
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.File(fileName)
}

func getCurrentStockReportAsExcel(c *gin.Context){

	//barcode := c.Query("barcode")
	name := c.Query("name")
	category := c.Query("category")

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")

	fileName, err := UseCase.GetCurrentStockReportAsExcel(name,category,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	println("file:",fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", `attachment; filename=excelFile.xlsx` )
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.File(fileName)
}


func getPaymentReportAsExcel (c *gin.Context){

	tInterval := c.Query("timeInterval")

	fileName, err := UseCase.GetPaymentReportAsExcel(tInterval)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	println("file:",fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", `attachment; filename=excelFile.xlsx` )
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.File(fileName)
}

func getProductReportAsExcel(c *gin.Context){

	tInterval := c.Query("timeInterval")
	category := c.Query("category")
	productName := c.Query("productName")
	userId,_ := strconv.Atoi(c.Query("userId"))

	fileName, err := UseCase.GetProductReportAsExcel(tInterval,productName,category,userId)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	println("file:",fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", `attachment; filename=excelFile.xlsx` )
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.File(fileName)
}


// ####################################################################

// # utils #

func handleLogin (c *gin.Context){
	p := LoginParams{}
	c.BindJSON(&p)

	u,t,err := UseCase.Login(p.Email,p.Password,secret)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(LoginResponse{Token: t, User: u}))

}

func generateSuccessResponse(data interface{}) (map[string]interface{}) {

	return gin.H{"data": data, "success": true, "errorCode": 0, "errorMessage": ""}
}

func generateFailResponse( err *ErrorType) (map[string]interface{}){
	return gin.H{"data": nil , "success": false, "errorCode": err.Code, "errorMessage": err.Message}
}