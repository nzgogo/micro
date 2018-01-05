package main

import (
	"log"
	micro "micro"
	"micro/router"
	"micro/codec"
	"micro/transport"
)

var(
	Codec = codec.NewCodec()
)

func ChangHenGe(req *codec.Request, tc transport.Transport, subject string) error {
	var rsp = "长恨歌 \n汉皇重色思倾国，御宇多年求不得。\n杨家有女初长成，养在深闺人未识。\n天生丽质难自弃，一朝选在君王侧。\n回眸一笑百媚生，六宫粉黛无颜色。\n春寒赐浴华清池，温泉水滑洗凝脂。\n侍儿扶起娇无力，始是新承恩泽时。\n云鬓花颜金步摇，芙蓉帐暖度春宵。\n春宵苦短日高起，从此君王不早朝。\n承欢侍宴无闲暇，春从春游夜专夜。\n后宫佳丽三千人，三千宠爱在一身。\n金屋妆成娇侍夜，玉楼宴罢醉和春。\n姊妹弟兄皆列土，可怜光彩生门户。\n骊宫高处入青云，仙乐风飘处处闻。\n缓歌谩舞凝丝竹，尽日君王看不足。\n渔阳鼙鼓动地来，惊破霓裳羽衣曲。\n翠华摇摇行复止，西出都门百余里。\n六军不发无奈何，宛转蛾眉马前死。\n花钿委地无人收，翠翘金雀玉搔头。\n君王掩面救不得，回看血泪相和流。\n黄埃散漫风萧索，云栈萦纡登剑阁。\n峨嵋山下少人行，旌旗无光日色薄。\n蜀江水碧蜀山青，圣主朝朝暮暮情。\n行宫见月伤心色，夜雨闻铃肠断声。\n天旋地转回龙驭，到此踌躇不能去。\n马嵬坡下泥土中，不见玉颜空死处。\n君臣相顾尽沾衣，东望都门信马归。\n归来池苑皆依旧，太液芙蓉未央柳。\n芙蓉如面柳如眉，对此如何不泪垂。\n春风桃李花开日，秋雨梧桐叶落时。\n西宫南内多秋草，落叶满阶红不扫。\n梨园弟子白发新，椒房阿监青娥老。\n夕殿萤飞思悄然，孤灯挑尽未成眠。\n迟迟钟鼓初长夜，耿耿星河欲曙天。\n鸳鸯瓦冷霜华重，翡翠衾寒谁与共。\n悠悠生死别经年，魂魄不曾来入梦。\n临邛道士鸿都客，能以精诚致魂魄。\n为感君王辗转思，遂教方士殷勤觅。\n排空驭气奔如电，升天入地求之遍。\n上穷碧落下黄泉，两处茫茫皆不见。\n忽闻海上有仙山，山在虚无缥渺间。\n楼阁玲珑五云起，其中绰约多仙子。\n中有一人字太真，雪肤花貌参差是。\n金阙西厢叩玉扃，转教小玉报双成。\n闻道汉家天子使，九华帐里梦魂惊。\n揽衣推枕起徘徊，珠箔银屏迤逦开。\n云鬓半偏新睡觉，花冠不整下堂来。\n风吹仙袂飘飘举，犹似霓裳羽衣舞。\n玉容寂寞泪阑干，梨花一枝春带雨。\n含情凝睇谢君王，一别音容两渺茫。\n昭阳殿里恩爱绝，蓬莱宫中日月长。\n回头下望人寰处，不见长安见尘雾。\n惟将旧物表深情，钿合金钗寄将去。\n钗留一股合一扇，钗擘黄金合分钿。\n但教心似金钿坚，天上人间会相见。\n临别殷勤重寄词，词中有誓两心知。\n七月七日长生殿，夜半无人私语时。\n在天愿作比翼鸟，在地愿为连理枝。\n天长地久有时尽，此恨绵绵无绝期。"
	response := &codec.Response{
		200,
		make(map[string][]string,0),
		rsp,
	}
	resp, err := Codec.Marshal(response)
	if err!=nil {
		return err
	}
	return tc.Publish(subject, resp)
}

func Hello(req *codec.Request, tc transport.Transport, subject string) error {
	response := &codec.Response{
		200,
		make(map[string][]string,0),
		"For the brave souls who get this far: You are the chosen ones, the valiant knights of programming who toil away, without rest, fixing our most awful code. To you, true saviors, kings of men, I say this: never gonna give you up, never gonna let you down, never gonna run around and desert you. Never gonna make you cry, never gonna say goodbye. Never gonna tell a lie and hurt you.",
	}
	resp, err := Codec.Marshal(response)
	if err!=nil {
		return err
	}
	return tc.Publish(subject, resp)
}

func main() {
	route := router.NewRouter(router.Name("gogox/v1/greeter"))
	route.Add(&router.Node{
		Method:"GET",
		Path:"/Changhenge",
		ID: "/Changhenge",
		Handler: ChangHenGe,
	})

	route.Add(&router.Node{
		Method:"GET",
		Path:"/hello",
		ID: "/hello",
		Handler: Hello,
	})

	service := micro.NewService(
		"gogox.core.greeter",
		"v1",
	)
	if err:=service.Init(micro.Router(route)); err!=nil{
		log.Fatal(err)
	}

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}