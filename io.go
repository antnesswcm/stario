package stario

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func Passwd(hint string, defaultVal string) InputMsg {
	return passwd(hint, defaultVal, "●")
}

func PasswdWithMask(hint string, defaultVal string, mask string) InputMsg {
	return passwd(hint, defaultVal, mask)
}

func MessageBoxRaw(hint string, defaultVal string) InputMsg {
	return messageBox(hint, defaultVal)
}

// 定义一个名为 messageBox 的函数，它接收两个字符串参数：hint 和 defaultVal，并返回一个 InputMsg 结构体
func messageBox(hint string, defaultVal string) InputMsg {
	var ioBuf []rune // 用于存储用户的输入字符序列
	if hint != "" {
		fmt.Print(hint) // 如果提示信息不为空，则输出提示信息
	}
	if strings.Index(hint, "\n") >= 0 {
		hint = strings.TrimSpace(hint[strings.LastIndex(hint, "\n"):])
	}
	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd) // 获取标准输入的原始状态，并将其设置为非标准模式
	if err != nil {
		return InputMsg{"", err} // 如果出现错误，则返回一个包含错误的结构体
	}
	defer fmt.Println()                      // 延迟执行 fmt.Println() 函数，确保最后输出一个换行符
	defer terminal.Restore(fd, state)        // 在函数返回之前恢复标准输入的状态
	inputReader := bufio.NewReader(os.Stdin) // 创建一个新的 bufio.Reader 对象，用于从标准输入读取数据
	for {
		b, _, err := inputReader.ReadRune() // 从标准输入读取一个 rune（字符）
		if err != nil {
			return InputMsg{"", err} // 如果出现错误，则返回一个包含错误的结构体
		}
		if b == 0x0d { // 如果读取的字符是回车符，则表示用户输入结束
			strValue := strings.TrimSpace(string(ioBuf)) // 将用户输入的字符序列转换为字符串，并去除字符串首尾的空白字符
			if len(strValue) == 0 {                      // 如果用户没有输入任何字符，则使用默认值
				strValue = defaultVal
			}
			return InputMsg{strValue, nil} // 返回一个包含输入字符串和空错误的结构体
		}
		if b == 0x08 || b == 0x7F { // 如果读取的字符是退格符或删除符，则表示用户要删除之前输入的字符
			if len(ioBuf) > 0 {
				ioBuf = ioBuf[:len(ioBuf)-1] // 从字符序列中删除最后一个字符
			}
			fmt.Print("\r") // 将光标移动到行首
			for i := 0; i < len(ioBuf)+2+len(hint); i++ {
				fmt.Print(" ") // 用空格清除屏幕上已输入的字符
			}
		} else {
			ioBuf = append(ioBuf, b)
		}
		fmt.Print("\r")
		if hint != "" {
			fmt.Print(hint)
		}
		fmt.Print(string(ioBuf))
	}
}

// passwd 函数用于从标准输入中读取密码，支持设置提示信息，缺省值和掩码。
// 输入参数：
// hint: string类型，提示信息字符串，用于向用户说明输入的目的和格式，可以为空字符串。
// defaultVal: string类型，缺省值，当用户不输入任何内容时，返回此缺省值。
// mask: string类型，掩码，用于在用户输入时掩盖真实的字符，可以为空字符串。
// 返回值：
// InputMsg: struct类型，包含两个字段，分别为用户输入的密码和错误信息，当无错误时错误信息为nil。
func passwd(hint string, defaultVal string, mask string) InputMsg {
	// 定义ioBuf变量用于存储用户输入的字符
	var ioBuf []rune
	// 如果有提示信息，则输出提示信息
	if hint != "" {
		fmt.Print(hint)
	}
	// 如果提示信息中包含换行符，则截取最后一行作为新的提示信息
	if strings.Index(hint, "\n") >= 0 {
		hint = strings.TrimSpace(hint[strings.LastIndex(hint, "\n"):])
	}
	// 获取标准输入的文件描述符，用于后续操作
	fd := int(os.Stdin.Fd())
	// 将标准输入设置为Raw模式，禁用输入缓冲和回显，确保每个输入字符都能立即读取
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return InputMsg{"", err}
	}
	defer fmt.Println()               // 输出换行符，以便后续输出不被输入字符覆盖
	defer terminal.Restore(fd, state) // 恢复标准输入的模式，以便后续输入能够正常工作
	// 创建一个bufio.Reader对象用于读取标准输入
	inputReader := bufio.NewReader(os.Stdin)
	for {
		// 从标准输入中读取一个字符
		b, _, err := inputReader.ReadRune()
		if err != nil {
			return InputMsg{"", err}
		}
		// 如果读取到回车符，则表示用户输入结束，返回用户输入的字符串
		if b == 0x0d {
			strValue := strings.TrimSpace(string(ioBuf))
			if len(strValue) == 0 {
				strValue = defaultVal
			}
			return InputMsg{strValue, nil}
		}
		// 如果读取到退格键或删除键，则删除ioBuf中最后一个字符，并在屏幕上清除该字符
		if b == 0x08 || b == 0x7F {
			if len(ioBuf) > 0 {
				ioBuf = ioBuf[:len(ioBuf)-1]
			}
			fmt.Print("\r") // 将光标移动到行首
			for i := 0; i < len(ioBuf)+2+len(hint); i++ {
				fmt.Print(" ") // 用空格清除屏幕上已输入的字符
			}
		} else {
			ioBuf = append(ioBuf, b)
		}
		fmt.Print("\r") // 将光标移动到行首
		if hint != "" {
			fmt.Print(hint)
		}
		for i := 0; i < len(ioBuf); i++ {
			fmt.Print(mask) // 在屏幕上用掩码替换已输入的字符
		}
	}
}

// MessageBox 函数用于在控制台中打印提示信息，等待用户输入，然后将输入的内容返回给调用者。
// 参数hint是要打印的提示信息，defaultVal是内容默认值。
func MessageBox(hint string, defaultVal string) InputMsg {
	// 如果有提示信息，则将其打印到控制台。
	if hint != "" {
		fmt.Print(hint)
	}

	// 创建一个bufio.Reader对象，用于从标准输入流(os.Stdin)中读取用户的输入。
	inputReader := bufio.NewReader(os.Stdin)

	// 读取用户输入的字符串，直到用户输入了换行符为止。如果读取过程中出现错误，则返回一个InputMsg对象，其中的err成员包含错误信息。
	str, err := inputReader.ReadString('\n')
	if err != nil {
		return InputMsg{"", err}
	}

	// 移除读取到的字符串中的空格。
	str = strings.TrimSpace(str)

	// 如果读取到的字符串长度为0，则使用默认值。
	if len(str) == 0 {
		str = defaultVal
	}

	// 返回一个InputMsg对象，其中的str成员包含用户输入的字符串，err成员为nil。
	return InputMsg{str, nil}
}

// GetYesNoInput 函数用于在控制台中打印提示信息，然后等待用户输入“Y”或“N”。
// 如果用户输入了Y返回true。如果用户输入N返回false
// 参数hint是要打印的提示信息，defaults是bool默认值。
func GetYesNoInput(hint string, defaults bool) bool {
	for {
		// 调用MessageBox函数获取用户输入的字符串并将其转换为大写形式。
		res := strings.ToUpper(MessageBox(hint, "").MustString())

		// 如果用户没有输入任何内容，则返回默认值。
		if res == "" {
			return defaults
		}

		// 截取用户输入字符串的第一个字符，并将其与“Y”和“N”进行比较。
		res = res[0:1]
		if res == "Y" {
			return true
		} else if res == "N" {
			return false
		}
	}
}

// WaitUntilString 函数会读取从标准输入流 os.Stdin 中输入的字符，直到满足 trigger 参数指定的字符串条件为止
// 如果 hint 参数为读取期间显示的提示信息
// 如果 repeat 参数为 true，则在读取期间，如果输入错误，将会重复显示提示信息
// 函数返回 error 类型，如果在读取期间发生错误，则会返回该错误
func WaitUntilString(hint string, trigger string, repeat bool) error {
	pressLen := len([]rune(trigger)) // 计算 trigger 字符串的长度，使用 rune 类型，防止一个字符被错误的计算成多个字符的长度
	if trigger == "" {
		pressLen = 1 // 如果 trigger 为空，则将字符长度设置为 1，因为读取任何一个字符都是满足停止条件
	}
	fd := int(os.Stdin.Fd()) // 获取标准输入流 os.Stdin 的文件描述符
	if hint != "" {
		fmt.Print(hint) // 如果 hint 不为空，则在读取期间显示提示信息
	}
	state, err := terminal.MakeRaw(fd) // 将终端的状态更改为 raw 模式，从而能够读取单个字符而不是行，这里的 state 变量是一个终端模式的结构体，它记录了终端的原始模式
	if err != nil {
		return err // 如果终端状态更改失败，则返回错误
	}
	defer terminal.Restore(fd, state)        // 无论何时，都需要将终端状态还原为原始状态
	inputReader := bufio.NewReader(os.Stdin) // 使用带缓存的读取器读取标准输入流
	i := 0                                   // 用于追踪当前已经读取的字符数
	for {
		b, _, err := inputReader.ReadRune() // 读取一个 rune 类型的字符
		if err != nil {
			return err // 如果读取错误，则返回该错误
		}
		if trigger == "" {
			break // 如果 trigger 为空，则立即停止读取
		}
		if b == []rune(trigger)[i] { // 如果读取到的字符与 trigger 字符串中对应位置的字符相同，则继续读取
			i++
			if i == pressLen { // 如果已经读取了 trigger 字符串的全部字符，则停止读取
				break
			}
			continue
		}
		i = 0                     // 如果读取到的字符与 trigger 字符串中对应位置的字符不同，则重置已经读取的字符数
		if hint != "" && repeat { // 如果 hint 不为空且 repeat 为 true，则在读取期间，如果输入错误，将会重复显示提示信息
			fmt.Print("\r\n")
			fmt.Print(hint)
		}
	}
	return nil // 读取结束，没有发生错误，返回 nil

}
