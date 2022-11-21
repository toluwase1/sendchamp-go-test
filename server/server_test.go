package server

//var (
//	server *Server
//)
//
//var testServer struct {
//	router  *gin.Engine
//	handler *Server
//}
//
//func TestMain(m *testing.M) {
//	if err := godotenv.Load(".env"); err != nil {
//		log.Printf("couldn't load env vars: %v", err)
//	}
//	fmt.Println("Starting server tests")
//	c, err := config.Load()
//	if err != nil {
//		log.Fatal(err)
//	}
//	testServer.handler = &Server{
//		Config: c,
//	}
//	testServer.handler.Config.JWTSecret = "testSecret"
//	testServer.router = testServer.handler.setupRouter()
//	exitCode := m.Run()
//	os.Exit(exitCode)
//}
//
//func AuthorizeTestUser(t *testing.T) (string, models.User) {
//	user, _ := randomUser(t)
//	user.IsEmailActive = true
//	accToken, err := jwt.GenerateToken(user.Email, testServer.handler.Config.JWTSecret)
//
//	require.NoError(t, err)
//	return accToken, user
//}
