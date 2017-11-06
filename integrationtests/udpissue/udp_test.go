package udpissue

import (
	"fmt"
	"net"
	"os/exec"

	_ "github.com/lucas-clemente/quic-clients" // download clients

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("UDP debugging", func() {
	for i := 0; i < 1000; i++ {
		run := i
		It(fmt.Sprintf("run %d", run), func() {
			done := make(chan struct{})

			addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
			Expect(err).ToNot(HaveOccurred())
			ln, err := net.ListenUDP("udp", addr)
			Expect(err).ToNot(HaveOccurred())

			go func() {
				defer GinkgoRecover()
				data := make([]byte, 100)
				_, _, err := ln.ReadFrom(data)
				Expect(err).ToNot(HaveOccurred())
				close(done)
			}()

			var session *Session
			sessionDone := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				command := exec.Command(
					clientPath,
					"--quic-version=39",
					"--host=127.0.0.1",
					fmt.Sprintf("--port=%d", ln.LocalAddr().(*net.UDPAddr).Port),
					"https://quic.clemente.io",
				)
				var err error
				session, err = Start(command, nil, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				close(sessionDone)
			}()
			<-done
			session.Kill()
			<-sessionDone
		})
	}
})
