package handler

import "time"

/*
AppChecker - Application Checker
*/
type AppChecker struct {
}

/*
HealthCheck - Check if the application is up and working
*/
func HealthCheck() string {
	a := AppChecker{}
	return "<html><body><h1>" + a.Desc() + "</h1></body></html>"
}

/*
Desc - HealthChecker App Description
*/
func (dc *AppChecker) Desc() string {
	t := time.Now()
	t.Format("2006-01-02 15:04:05")
	return "[1] [mercurius] Tudo Certo em " + t.String()
}

/*
Check - HealthChecker check App status function
*/
func (dc *AppChecker) Check() error {
	return nil
}