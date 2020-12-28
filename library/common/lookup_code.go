package common

import (
	"strings"

	"github.com/gogf/gf/container/gmap"
)


var (
	Countrys *gmap.StrAnyMap
	data = `Angola,安哥拉,AO,244,-7
Afghanistan,阿富汗,AF,93,0
Albania,阿尔巴尼亚,AL,355,-7
Algeria,阿尔及利亚,DZ,213,-8
Andorra,安道尔共和国,AD,376,-8
Anguilla,安圭拉岛,AI,1264,-12
Antigua and Barbuda,安提瓜和巴布达,AG,1268,-12
Argentina,阿根廷,AR,54,-11
Armenia,亚美尼亚,AM,374,-6
Ascension,阿森松, ,247,-8
Australia,澳大利亚,AU,61,+2
Austria,奥地利,AT,43,-7
Azerbaijan,阿塞拜疆,AZ,994,-5
Bahamas,巴哈马,BS,1242,-13
Bahrain,巴林,BH,973,-5
Bangladesh,孟加拉国,BD,880,-2
Barbados,巴巴多斯,BB,1246,-12
Belarus,白俄罗斯,BY,375,-6
Belgium,比利时,BE,32,-7
Belize,伯利兹,BZ,501,-14
Benin,贝宁,BJ,229,-7
Bermuda Is.,百慕大群岛,BM,1441,-12
Bolivia,玻利维亚,BO,591,-12
Botswana,博茨瓦纳,BW,267,-6
Brazil,巴西,BR,55,-11
Brunei,文莱,BN,673,0
Bulgaria,保加利亚,BG,359,-6
Burkina-faso,布基纳法索,BF,226,-8
Burma,缅甸,MM,95,-1.3
Burundi,布隆迪,BI,257,-6
Cameroon,喀麦隆,CM,237,-7
Canada,加拿大,CA,1,-13
Cayman Is.,开曼群岛, ,1345,-13
Central African Republic,中非共和国,CF,236,-7
Chad,乍得,TD,235,-7
Chile,智利,CL,56,-13
China,中国,CN,86,0
Colombia,哥伦比亚,CO,57,0
Congo,刚果,CG,242,-7
Cook Is.,库克群岛,CK,682,-18.3
Costa Rica,哥斯达黎加,CR,506,-14
Cuba,古巴,CU,53,-13
Cyprus,塞浦路斯,CY,357,-6
Czech Republic,捷克,CZ,420,-7
Denmark,丹麦,DK,45,-7
Djibouti,吉布提,DJ,253,-5
Dominica Rep.,多米尼加共和国,DO,1890,-13
Ecuador,厄瓜多尔,EC,593,-13
Egypt,埃及,EG,20,-6
EI Salvador,萨尔瓦多,SV,503,-14
Estonia,爱沙尼亚,EE,372,-5
Ethiopia,埃塞俄比亚,ET,251,-5
Fiji,斐济,FJ,679,+4
Finland,芬兰,FI,358,-6
France,法国,FR,33,-8
French Guiana,法属圭亚那,GF,594,-12
Gabon,加蓬,GA,241,-7
Gambia,冈比亚,GM,220,-8
Georgia,格鲁吉亚,GE,995,0
Germany,德国,DE,49,-7
Ghana,加纳,GH,233,-8
Gibraltar,直布罗陀,GI,350,-8
Greece,希腊,GR,30,-6
Grenada,格林纳达,GD,1809,-14
Guam,关岛,GU,1671,+2
Guatemala,危地马拉,GT,502,-14
Guinea,几内亚,GN,224,-8
Guyana,圭亚那,GY,592,-11
Haiti,海地,HT,509,-13
Honduras,洪都拉斯,HN,504,-14
Hongkong,香港,HK,852,0
Hungary,匈牙利,HU,36,-7
Iceland,冰岛,IS,354,-9
India,印度,IN,91,-2.3
Indonesia,印度尼西亚,ID,62,-0.3
Iran,伊朗,IR,98,-4.3
Iraq,伊拉克,IQ,964,-5
Ireland,爱尔兰,IE,353,-4.3
Israel,以色列,IL,972,-6
Italy,意大利,IT,39,-7
Ivory Coast,科特迪瓦, ,225,-6
Jamaica,牙买加,JM,1876,-12
Japan,日本,JP,81,+1
Jordan,约旦,JO,962,-6
Kampuchea (Cambodia ),柬埔寨,KH,855,-1
Kazakstan,哈萨克斯坦,KZ,327,-5
Kenya,肯尼亚,KE,254,-5
Korea,韩国,KR,82,+1
Kuwait,科威特,KW,965,-5
Kyrgyzstan,吉尔吉斯坦,KG,331,-5
Laos,老挝,LA,856,-1
Latvia,拉脱维亚,LV,371,-5
Lebanon,黎巴嫩,LB,961,-6
Lesotho,莱索托,LS,266,-6
Liberia,利比里亚,LR,231,-8
Libya,利比亚,LY,218,-6
Liechtenstein,列支敦士登,LI,423,-7
Lithuania,立陶宛,LT,370,-5
Luxembourg,卢森堡,LU,352,-7
Macao,澳门,MO,853,0
Madagascar,马达加斯加,MG,261,-5
Malawi,马拉维,MW,265,-6
Malaysia,马来西亚,MY,60,-0.5
Maldives,马尔代夫,MV,960,-7
Mali,马里,ML,223,-8
Malta,马耳他,MT,356,-7
Mariana Is,马里亚那群岛, ,1670,+1
Martinique,马提尼克, ,596,-12
Mauritius,毛里求斯,MU,230,-4
Mexico,墨西哥,MX,52,-15
Moldova Republic of,摩尔多瓦,MD,373,-5
Monaco,摩纳哥,MC,377,-7
Mongolia,蒙古,MN,976,0
Montserrat Is,蒙特塞拉特岛,MS,1664,-12
Morocco,摩洛哥,MA,212,-6
Mozambique,莫桑比克,MZ,258,-6
Namibia,纳米比亚,NA,264,-7
Nauru,瑙鲁,NR,674,+4
Nepal,尼泊尔,NP,977,-2.3
Netheriands Antilles,荷属安的列斯, ,599,-12
Netherlands,荷兰,NL,31,-7
New Zealand,新西兰,NZ,64,+4
Nicaragua,尼加拉瓜,NI,505,-14
Niger,尼日尔,NE,227,-8
Nigeria,尼日利亚,NG,234,-7
North Korea,朝鲜,KP,850,+1
Norway,挪威,NO,47,-7
Oman,阿曼,OM,968,-4
Pakistan,巴基斯坦,PK,92,-2.3
Panama,巴拿马,PA,507,-13
Papua New Cuinea,巴布亚新几内亚,PG,675,+2
Paraguay,巴拉圭,PY,595,-12
Peru,秘鲁,PE,51,-13
Philippines,菲律宾,PH,63,0
Poland,波兰,PL,48,-7
French Polynesia,法属玻利尼西亚,PF,689,+3
Portugal,葡萄牙,PT,351,-8
Puerto Rico,波多黎各,PR,1787,-12
Qatar,卡塔尔,QA,974,-5
Reunion,留尼旺, ,262,-4
Romania,罗马尼亚,RO,40,-6
Russia,俄罗斯,RU,7,-5
Saint Lueia,圣卢西亚,LC,1758,-12
Saint Vincent,圣文森特岛,VC,1784,-12
Samoa Eastern,东萨摩亚(美), ,684,-19
Samoa Western,西萨摩亚, ,685,-19
San Marino,圣马力诺,SM,378,-7
Sao Tome and Principe,圣多美和普林西比,ST,239,-8
Saudi Arabia,沙特阿拉伯,SA,966,-5
Senegal,塞内加尔,SN,221,-8
Seychelles,塞舌尔,SC,248,-4
Sierra Leone,塞拉利昂,SL,232,-8
Singapore,新加坡,SG,65,+0.3
Slovakia,斯洛伐克,SK,421,-7
Slovenia,斯洛文尼亚,SI,386,-7
Solomon Is,所罗门群岛,SB,677,+3
Somali,索马里,SO,252,-5
South Africa,南非,ZA,27,-6
Spain,西班牙,ES,34,-8
Sri Lanka,斯里兰卡,LK,94,0
St.Lucia,圣卢西亚,LC,1758,-12
St.Vincent,圣文森特,VC,1784,-12
Sudan,苏丹,SD,249,-6
Suriname,苏里南,SR,597,-11.3
Swaziland,斯威士兰,SZ,268,-6
Sweden,瑞典,SE,46,-7
Switzerland,瑞士,CH,41,-7
Syria,叙利亚,SY,963,-6
Taiwan,台湾省,TW,886,0
Tajikstan,塔吉克斯坦,TJ,992,-5
Tanzania,坦桑尼亚,TZ,255,-5
Thailand,泰国,TH,66,-1
Togo,多哥,TG,228,-8
Tonga,汤加,TO,676,+4
Trinidad and Tobago,特立尼达和多巴哥,TT,1809,-12
Tunisia,突尼斯,TN,216,-7
Turkey,土耳其,TR,90,-6
Turkmenistan,土库曼斯坦,TM,993,-5
Uganda,乌干达,UG,256,-5
Ukraine,乌克兰,UA,380,-5
United Arab Emirates,阿拉伯联合酋长国,AE,971,-4
United Kiongdom,英国,GB,44,-8
United States of America,美国,US,1,-13
Uruguay,乌拉圭,UY,598,-10.3
Uzbekistan,乌兹别克斯坦,UZ,233,-5
Venezuela,委内瑞拉,VE,58,-12.3
Vietnam,越南,VN,84,-1
Yemen,也门,YE,967,-5
Yugoslavia,南斯拉夫,YU,381,-7
Zimbabwe,津巴布韦,ZW,263,-6
Zaire,扎伊尔,ZR,243,-7
Zambia,赞比亚,ZM,260,-6`
)

func init() {
	Countrys = gmap.NewStrAnyMap()

	for _, s := range strings.Split(data, "\n") {
		i := strings.Split(s, ",")
		c := &Country{
			EN: i[0],
			CN: i[1],
			Code: i[2],
			TimeZone: i[4],
		}
		Countrys.Set(c.Code, c)
	}
}

type Country struct {
	EN       string `json:"en"`
	CN       string `json:"cn"`
	Code     string `json:"code"`
	TimeZone string `json:"time_zone"`
}

func SearchCountryFromCode(code string) *Country {
	if Countrys.Contains(code) {
		v := Countrys.Get(code)
		return v.(*Country)
	} else {
		return nil
	}
}