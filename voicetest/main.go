// voicetest project main.go
package main

import (
	"baiduai"
	"fmt"
	"ghostlib"
	"io/ioutil"
)

/**
 * 语音合成
 */
func test_text2voice() {
	voice := baiduai.NewVoice()
	T := `
晚上跑完步有点口渴，到路边摊买橘子，挑了几个卖相好的。
老板又拿起几个有斑点的说：“帅哥，这种长得不好看的其实更甜。”
我颇有感悟的说：“是因为橘子觉得自己长的不好看，所以努力让自己变得更甜吗？”
老板微微一愣道：“不是，我想早点卖完回家。”
`
	flag, result := voice.GetVoice(T)
	if !flag {
		ghostlib.Msg(string(result), 3)
	} else {
		fmt.Printf("\nlen:%v", len(result))
		err := ioutil.WriteFile("E:/data/bd_v.wav", result, 0766)
		if err != nil {
			ghostlib.Msg("写入结果文件[/data1/bd_v.mp3]出错", 3)
		}
	}
}

/**
 * 语音识别
 */
func test_voice2txt() {
	engine := baiduai.NewVoice()
	// bd_voice.wav 输入的是 单声道16k采样
	r, err := ioutil.ReadFile("E:/data/bd_v.wav")
	if nil != err {
		panic(err)
	}

	txtresult, err := engine.GetText(r)
	if nil != err {
		fmt.Printf("\n%v\n", err)
	}
	fmt.Printf("%v\n", txtresult)
}

/**
 * 评论观点抽取 或 情感识别  垃圾 很容易没结果
 */
func test_ctag() {
	T := "京东商城的商品质量很好,价格也非常便宜,京东商城的快递很迅速,快递员服务态度也非常好,在京东商城购物省时省力省钱"
	T = "辛苦快递小哥了！抢购的！相信京东没喝应该是正品！个人觉得盒子好小，毕竟只有400Ml幸好有纸质包装袋，可以装一下，才敢送出手！不管怎么说这个牌子广告响！值吧！就是赠品一个都没送，还买了3箱"
	engine := baiduai.NewText()
	x := engine.GetCommentTag(T)
	//fmt.Println(x)

	fmt.Printf("\n---------\nsrc:%s\n------\n", T)
	for _, r := range x {
		r := r.(map[string]interface{})
		fmt.Printf("abstract: %s\n", r["abstract"].(string))
		fmt.Printf("fea: %s\n", r["fea"].(string))
		fmt.Printf("adj: %s\n-------\n", r["adj"].(string))
	}
}

/**
 * 分词
 */
func test_splitword() {
	T := `
王辉。中国美术家协会会员。1980年生于湖北。2003年毕业于湖北美术学院。
中国画专业。学士学位。2006年毕业于广西艺术学院。中国画专业。硕士学位。
现任职于四川绵阳师范学院美术学院。中国画教师。作品多次入选全国性学术展览并获奖。
    `
	engine := baiduai.NewText()
	x := engine.SplitWords(T)

	namebuf := x["namebuf"]                                            // 人名
	subphrbuf := x["subphrbuf"]                                        // 短语
	wordsepbuf := x["wordsepbuf"]                                      // 标准粒度
	wpcompbuf := x["wpcompbuf"]                                        // 混排粒度
	newwordbuf := x["pnewword"].(map[string]interface{})["newwordbuf"] // 新词

	fmt.Printf("src:%v", T)
	fmt.Printf("\n人名: %v", namebuf)
	fmt.Printf("\n短语: %v", subphrbuf)
	fmt.Printf("\n新词: %v", newwordbuf)
	fmt.Printf("\n标准: %v", wordsepbuf)
	fmt.Printf("\n混排: %v", wpcompbuf)

}

/**
 * 词性标注
 */
func test_wordpos() {
	T := `
王辉。中国美术家协会会员。1980年生于湖北。2003年毕业于湖北美术学院。
中国画专业。学士学位。2006年毕业于广西艺术学院。中国画专业。硕士学位。
现任职于四川绵阳师范学院美术学院。中国画教师。作品多次入选全国性学术展览并获奖。
    `
	engine := baiduai.NewText()
	x := engine.WordPos(T)
	for _, r := range x {
		r := r.(map[string]interface{})
		// 非符号的词
		if "w" != r["type"].(string) {
			fmt.Printf("word: %s\n", r["word"].(string))
			fmt.Printf("type: %s\n", r["type"].(string))
			fmt.Printf("kind: %s\n-------\n", engine.GetWordKind(r["type"].(string)))

		}
	}
}

/**
 * dnn中文语言模型
 */
func test_dnnlm() {
	T := `
王辉。中国美术家协会会员。1980年生于湖北。2003年毕业于湖北美术学院。
中国画专业。学士学位。2006年毕业于广西艺术学院。中国画专业。硕士学位。
现任职于四川绵阳师范学院美术学院。中国画教师。作品多次入选全国性学术展览并获奖。
    `
	engine := baiduai.NewText()
	x := engine.DnnLm(T)
	fmt.Printf("%v", x)
}

/**
 * 短文本相似度
 */
func test_simnet() {
	a := "我是你爸爸"
	b := "你是我爸爸也不行"
	engine := baiduai.NewText()
	x := engine.SimNet(a, b)
	fmt.Printf("%v", x)
}

/**
 * 身份证识别
 */
func test_ocridcard() {
	engine := baiduai.NewOcr()
	r, err := ioutil.ReadFile("E:/data/s1.jpg")
	if nil != err {
		panic(err)
	}
	x := engine.OcrIdCard(r, true)
	//fmt.Printf("%v\n", x)
	result := x["words_result"].(map[string]interface{})
	for kk, vv := range result {
		fmt.Printf("--------\n%v:%v\n", kk, vv.(map[string]interface{})["words"])
	}
}

/**
 * 银行卡识别
 */
func test_ocrbankcard() {
	engine := baiduai.NewOcr()
	r, err := ioutil.ReadFile("E:/data/b2.jpg")
	if nil != err {
		panic(err)
	}
	x := engine.OcrBankCard(r)
	fmt.Printf("%v\n", x)
	result := x["result"].(map[string]interface{})
	for kk, vv := range result {
		fmt.Printf("--------\n%v:%v\n", kk, vv)
	}
}

/**
 * 文字识别
 */
func test_ocrgeneral() {
	engine := baiduai.NewOcr()
	r, err := ioutil.ReadFile("E:/data/11.jpg")
	if nil != err {
		panic(err)
	}
	x := engine.OcrGeneral(r)
	//fmt.Printf("%v\n", x)
	result := x["words_result"].([]interface{})
	for kk, vv := range result {
		fmt.Printf("--------\n%v:%v\n", kk, vv.(map[string]interface{})["words"])
	}
}

/**
 * 人脸检测
 */
func test_facedetect() {
	engine := baiduai.NewFace()
	r, err := ioutil.ReadFile("E:/data/f1.jpg")
	if nil != err {
		panic(err)
	}
	x := engine.FaceDetect(r)
	fmt.Printf("%v\n", x)

}

/**
 * 人脸相似度匹配
 */
func test_facematch() {
	engine := baiduai.NewFace()
	r, err := ioutil.ReadFile("E:/data/f2.jpg")
	if nil != err {
		panic(err)
	}
	r1, err := ioutil.ReadFile("E:/data/f3.jpg")
	if nil != err {
		panic(err)
	}
	r2, err := ioutil.ReadFile("E:/data/f4.jpg")
	if nil != err {
		panic(err)
	}
	r3, err := ioutil.ReadFile("E:/data/f5.jpg")
	if nil != err {
		panic(err)
	}
	x := engine.FaceMatch(r, r1, r2, r3)
	//fmt.Printf("%v\n", x)

	results := x["results"].([]interface{})
	for i, r := range results {
		r := r.(map[string]interface{})
		fmt.Printf("\n--------第%v组比对----------\n", i)
		index_i := ghostlib.ToInt64(r["index_i"]) + 1
		index_j := ghostlib.ToInt64(r["index_j"]) + 1
		fmt.Printf("img[%v]与img[%v] 相似度:%v%%\n", index_i, index_j, r["score"])
	}
}

/**
 * 色情识别
 */
func test_antporn() {
	engine := baiduai.NewFace()
	//r, err := ioutil.ReadFile("/data1/s1.jpg")
	r, err := ioutil.ReadFile("/data1/s9.jpg")
	if nil != err {
		panic(err)
	}
	x := engine.AntiPorn(r)

	var max_class = ""
	var max_prob float64 = 0

	for _, r := range x["result"].([]interface{}) {
		r := r.(map[string]interface{})
		cprob := r["probability"].(float64)
		if cprob > max_prob {
			max_prob = cprob
			max_class = r["class_name"].(string)
		}

		fmt.Printf("\n-------------\n分类: %v\n", r["class_name"])
		fmt.Printf("置信度: %v\n", r["probability"])
	}

	fmt.Printf("\n********结果*******\n")
	fmt.Printf("这是 [ %v ] 图片的可能性为 [ %v%% ]", max_class, max_prob*100)

}

func main() {
	//文字转语音测试
	//test_text2voice()
	//语音转文字测试
	//test_voice2txt()
	//test_ctag()
	//test_splitword()
	//test_wordpos()
	//test_dnnlm()
	//test_simnet()
	//身份证识别测试(照片必须是正面照)
	//test_ocridcard()
	//银行卡识别测试
	test_ocrbankcard()
	//test_ocrgeneral()

	//test_facedetect()
	//test_facematch()

	//test_antporn()
}
