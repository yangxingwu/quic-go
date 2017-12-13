package wire

import (
	"bytes"
	"io"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/qerr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("APPLICATION_CLOSE Frame", func() {
	Context("when parsing", func() {
		It("accepts sample frame", func() {
			data := []byte{0x3, 0x0, 0x19}
			data = append(data, encodeVarInt(0x1b)...) // reason phrase length
			data = append(data, []byte{
				'N', 'o', ' ', 'r', 'e', 'c', 'e', 'n', 't', ' ', 'n', 'e', 't', 'w', 'o', 'r', 'k', ' ', 'a', 'c', 't', 'i', 'v', 'i', 't', 'y', '.',
			}...)
			b := bytes.NewReader(data)
			frame, err := ParseApplicationCloseFrame(b, versionIETFFrames)
			Expect(err).ToNot(HaveOccurred())
			Expect(frame.ErrorCode).To(Equal(qerr.ErrorCode(0x19)))
			Expect(frame.ReasonPhrase).To(Equal("No recent network activity."))
			Expect(b.Len()).To(BeZero())
		})

		It("rejects long reason phrases", func() {
			data := []byte{0x3, 0xca, 0xfe}
			data = append(data, encodeVarInt(0xffff)...) // reason phrase length
			b := bytes.NewReader(data)
			_, err := ParseApplicationCloseFrame(b, versionIETFFrames)
			Expect(err).To(MatchError(io.EOF))
		})

		It("errors on EOFs", func() {
			data := []byte{0x3, 0x0, 0x19}
			data = append(data, encodeVarInt(0x1b)...) // reason phrase length
			data = append(data, []byte{
				'N', 'o', ' ', 'r', 'e', 'c', 'e', 'n', 't', ' ', 'n', 'e', 't', 'w', 'o', 'r', 'k', ' ', 'a', 'c', 't', 'i', 'v', 'i', 't', 'y', '.',
			}...)
			_, err := ParseApplicationCloseFrame(bytes.NewReader(data), versionIETFFrames)
			Expect(err).NotTo(HaveOccurred())
			for i := range data {
				_, err := ParseApplicationCloseFrame(bytes.NewReader(data[0:i]), versionIETFFrames)
				Expect(err).To(HaveOccurred())
			}
		})

		It("parses a frame without a reason phrase", func() {
			data := []byte{0x3, 0xca, 0xfe}
			data = append(data, encodeVarInt(0)...)
			b := bytes.NewReader(data)
			frame, err := ParseApplicationCloseFrame(b, versionIETFFrames)
			Expect(err).ToNot(HaveOccurred())
			Expect(frame.ReasonPhrase).To(BeEmpty())
			Expect(b.Len()).To(BeZero())
		})
	})

	Context("when writing", func() {
		It("writes a frame without a ReasonPhrase", func() {
			b := &bytes.Buffer{}
			frame := &ApplicationCloseFrame{
				ErrorCode: 0xbeef,
			}
			err := frame.Write(b, versionIETFFrames)
			Expect(err).ToNot(HaveOccurred())
			expected := []byte{0x3, 0xbe, 0xef}
			expected = append(expected, encodeVarInt(0)...)
			Expect(b.Bytes()).To(Equal(expected))
		})

		It("writes a frame with a ReasonPhrase", func() {
			b := &bytes.Buffer{}
			frame := &ApplicationCloseFrame{
				ErrorCode:    0xdead,
				ReasonPhrase: "foobar",
			}
			err := frame.Write(b, versionIETFFrames)
			Expect(err).ToNot(HaveOccurred())
			expected := []byte{0x3, 0xde, 0xad}
			expected = append(expected, encodeVarInt(6)...)
			expected = append(expected, []byte{'f', 'o', 'o', 'b', 'a', 'r'}...)
			Expect(b.Bytes()).To(Equal(expected))
		})

		It("has proper min length", func() {
			b := &bytes.Buffer{}
			f := &ApplicationCloseFrame{
				ErrorCode:    0xcafe,
				ReasonPhrase: "foobar",
			}
			err := f.Write(b, versionIETFFrames)
			Expect(err).ToNot(HaveOccurred())
			Expect(f.MinLength(versionIETFFrames)).To(Equal(protocol.ByteCount(b.Len())))
		})
	})
})
