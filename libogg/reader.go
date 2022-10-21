package libogg

/*
#cgo pkg-config: ogg
#include "ogg/ogg.h"
*/
import "C"
import (
	"bufio"
	"fmt"
	"unsafe"
)

type OggReader struct {
	r            *bufio.Reader
	sync_state   C.ogg_sync_state
	stream_state C.ogg_stream_state
	page         C.ogg_page
}

func NewOggReader(r *bufio.Reader) (*OggReader, error) {
	reader := OggReader{
		r: r,
	}

	C.ogg_sync_init(&reader.sync_state)

	return &reader, nil
}

func (r *OggReader) Next() bool {

	ret := C.ogg_stream_packetpeek(&r.stream_state, nil)
	if ret == 1 {
		return true
	}

	ret = C.ogg_sync_pageout(&r.sync_state, &r.page)
	if ret != 1 {
		for {
			buff := C.ogg_sync_buffer(&r.sync_state, 8192)
			scanData := (*[8192]byte)(unsafe.Pointer(buff))[:8192]

			nRead, err := r.r.Read(scanData)
			if err != nil {
				//return err
				return false
			}
			fmt.Printf("read: %d\n", nRead)
			ret = C.ogg_sync_wrote(&r.sync_state, C.long(nRead))
			if ret == -1 {
				fmt.Printf("sync wrote res: %d", ret)
				return false
			}

			ret = C.ogg_sync_pageout(&r.sync_state, &r.page)
			if ret != 1 {
				fmt.Printf("sync page out: %d", ret)
				//return false
				continue
			}
			fmt.Printf("page %d has packets: %d\n", C.ogg_page_pageno(&r.page), C.ogg_page_packets(&r.page))

			if int(C.ogg_page_pageno(&r.page)) == 0 {
				ret = C.ogg_stream_init(&r.stream_state, C.ogg_page_serialno(&r.page))
				if ret != 0 {
					fmt.Printf("stream init: %d", ret)
					return false
				}
			}
			break
		}
	}

	ret = C.ogg_stream_pagein(&r.stream_state, &r.page)
	if ret != 0 {
		fmt.Printf("stream page in: %d", ret)
		return false
	}

	ret = C.ogg_stream_packetpeek(&r.stream_state, nil)

	return ret == 1
}

func (r *OggReader) Scan() ([]byte, error) {
	var pkt C.ogg_packet

	ret := C.ogg_stream_packetout(&r.stream_state, &pkt)
	if ret != 1 {
		return nil, fmt.Errorf("stream packet out: %d", ret)
	}

	pktData := (*[8192]byte)(unsafe.Pointer(pkt.packet))[:int(pkt.bytes)]

	return pktData, nil
}

func (r *OggReader) Close() {
	C.ogg_sync_clear(&r.sync_state)
}
