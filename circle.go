package stario

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

// 定义一个名为 StarBuffer 的结构体类型，该类型实现了 io.Reader, io.Writer 和 io.Closer 接口
type StarBuffer struct {
	io.Reader            // 读接口
	io.Writer            // 写接口
	io.Closer            // 关闭接口
	datas     []byte     // 存储数据的字节数组
	pStart    int        // 缓冲区的起始位置
	pEnd      int        // 缓冲区的结束位置
	cap       int        // 缓冲区的容量
	isClose   bool       // 缓冲区是否关闭
	isEnd     bool       // 缓冲区是否已满
	rmu       sync.Mutex // 读互斥锁
	wmu       sync.Mutex // 写互斥锁
}

// NewStarBuffer 是一个构造函数，创建一个指定容量的 StarBuffer 对象
func NewStarBuffer(cap int) *StarBuffer {
	rtnBuffer := new(StarBuffer)
	rtnBuffer.cap = cap
	rtnBuffer.datas = make([]byte, cap)
	return rtnBuffer
}

// Free 方法返回缓冲区中可用空间的大小
func (star *StarBuffer) Free() int {
	return star.cap - star.Len()
}

// Cap 方法返回缓冲区的容量
func (star *StarBuffer) Cap() int {
	return star.cap
}

// Len 方法返回缓冲区中当前的数据长度
func (star *StarBuffer) Len() int {
	length := star.pEnd - star.pStart
	if length < 0 { // 缓冲区已经循环
		return star.cap + length - 1
	}
	return length
}

// getByte 方法从缓冲区中读取一个字节并返回，如果缓冲区已关闭或已满并且当前数据长度为0，则返回 io.EOF 错误。
// 如果当前数据长度为0，则返回 errors.New("no byte available now") 错误。
func (star *StarBuffer) getByte() (byte, error) {
	if star.isClose || (star.isEnd && star.Len() == 0) {
		return 0, io.EOF
	}
	if star.Len() == 0 {
		return 0, errors.New("no byte available now")
	}
	data := star.datas[star.pStart] // 读取数据
	star.pStart++                   // 指向下一个字节
	if star.pStart == star.cap {    // 如果已经到缓冲区的末尾，则循环到缓冲区的开头
		star.pStart = 0
	}
	return data, nil // 返回读取的数据
}

// putByte 方法将一个字节写入缓冲区，如果缓冲区已关闭或已满，则返回 io.EOF 错误。
// 如果缓冲区中没有可用空间，则会一直等待，直到有空间可用。
func (star *StarBuffer) putByte(data byte) error {
	if star.isClose || star.isEnd {
		return io.EOF
	}
	kariEnd := star.pEnd + 1 // 计算下一个可写入数据的位置
	if kariEnd == star.cap { // 如果已经到缓冲区的末尾，则循环到缓冲区的开头
		kariEnd = 0
	}
	if kariEnd == star.pStart { // 如果没有可用空间，则一直等待
		for {
			time.Sleep(time.Microsecond)
			if kariEnd != star.pStart {
				break
			}
		}
	}
	star.datas[star.pEnd] = data // 写入数据
	star.pEnd = kariEnd          // 更新下一个可写入数据的位置
	return nil
}

// Close 方法关闭缓冲区
func (star *StarBuffer) Close() error {
	star.isClose = true // 标记缓冲区已关闭
	return nil
}

// Read 从StarBuffer中读取数据到buf中，返回读取的字节数和错误信息
func (star *StarBuffer) Read(buf []byte) (int, error) {
	if star.isClose { // 如果StarBuffer已经关闭，则返回io.EOF
		return 0, io.EOF
	}
	if buf == nil { // 如果buf为空，则返回错误信息
		return 0, errors.New("buffer is nil")
	}
	star.rmu.Lock()                 // 获取读锁
	defer star.rmu.Unlock()         // 释放读锁
	var sum int = 0                 // 定义变量sum来计算读取的字节数
	for i := 0; i < len(buf); i++ { // 循环读取buf中的数据
		data, err := star.getByte() // 从StarBuffer中获取一个字节
		if err != nil {             // 如果获取字节时出错
			if err == io.EOF { // 如果错误是EOF，则返回sum和错误信息
				return sum, err
			}
			return sum, nil // 否则只返回sum
		}
		buf[i] = data // 将获取的字节写入buf中
		sum++         // sum加1
	}
	return sum, nil // 返回读取的字节数和错误信息（如果没有错误则为nil）
}

// Write 将bts写入StarBuffer，返回写入的字节数和错误信息
func (star *StarBuffer) Write(bts []byte) (int, error) {
	if bts == nil || star.isClose { // 如果bts为空或StarBuffer已经关闭，则将isEnd设置为true，并返回io.EOF
		star.isEnd = true
		return 0, io.EOF
	}
	star.wmu.Lock()                 // 获取写锁
	defer star.wmu.Unlock()         // 释放写锁
	var sum = 0                     // 定义变量sum来计算写入的字节数
	for i := 0; i < len(bts); i++ { // 循环写入bts中的数据
		err := star.putByte(bts[i]) // 向StarBuffer中写入一个字节
		if err != nil {             // 如果写入时出错
			fmt.Println("Write bts err:", err)
			return sum, err // 返回sum和错误信息
		}
		sum++ // sum加1
	}
	return sum, nil // 返回写入的字节数和错误信息（如果没有错误则为nil）
}
