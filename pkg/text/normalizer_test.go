package text

import (
	"reflect"
	"testing"
)

func Test_normalizer_push(t *testing.T) {
	type args struct {
		p *JapanZipCode
	}
	tests := []struct {
		name   string
		normer *normalizer
		args   args
	}{
		{"no error", newNormalizer(), args{&JapanZipCode{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.normer.push(tt.args.p)
		})
	}
}

func Test_normalizer_canPop(t *testing.T) {

	n1 := newNormalizer()
	n2 := newNormalizer()
	n2.outputs = append(n2.outputs, &JapanZipCode{})

	tests := []struct {
		name   string
		normer *normalizer
		want   bool
	}{
		{"false", n1, false},
		{"true", n2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.normer.canPop(); got != tt.want {
				t.Errorf("normalizer.canPop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizer_pop(t *testing.T) {
	z1 := &JapanZipCode{PrefCode: "01"}
	z2 := &JapanZipCode{PrefCode: "02"}
	n1 := newNormalizer()
	n2 := newNormalizer()
	n2.outputs = append(n2.outputs, z1)
	n2.outputs = append(n2.outputs, z2)

	tests := []struct {
		name   string
		normer *normalizer
		want   *JapanZipCode
	}{
		{"", n1, nil},
		{"", n2, z1},
		{"", n2, z2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.normer.pop(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("normalizer.pop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizer_normalize(t *testing.T) {

	beforeAfter := [][][]string{
		{
			{`01101,"060  ","0600000","ﾎｯｶｲﾄﾞｳ","ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ","ｲｶﾆｹｲｻｲｶﾞﾅｲﾊﾞｱｲ","北海道","札幌市中央区","以下に掲載がない場合",0,0,0,0,0,0`},
			{`01101,"060  ","0600000","ﾎｯｶｲﾄﾞｳ","ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ","","北海道","札幌市中央区","",0,0,0,0,0,0`},
		},
		{
			{`08546,"30604","3060433","ｲﾊﾞﾗｷｹﾝ","ｻｼﾏｸﾞﾝｻｶｲﾏﾁ","ｻｶｲﾏﾁﾉﾂｷﾞﾆﾊﾞﾝﾁｶﾞｸﾙﾊﾞｱｲ","茨城県","猿島郡境町","境町の次に番地がくる場合",0,0,0,0,0,0`},
			{`08546,"30604","3060433","ｲﾊﾞﾗｷｹﾝ","ｻｼﾏｸﾞﾝｻｶｲﾏﾁ","","茨城県","猿島郡境町","",0,0,0,0,0,0`},
		},
		{
			{`13362,"10003","1000301","ﾄｳｷｮｳﾄ","ﾄｼﾏﾑﾗ","ﾄｼﾏﾑﾗｲﾁｴﾝ","東京都","利島村","利島村一円",0,0,0,0,0,0`},
			{`13362,"10003","1000301","ﾄｳｷｮｳﾄ","ﾄｼﾏﾑﾗ","","東京都","利島村","",0,0,0,0,0,0`},
		},
		{
			{`25443,"52203","5220317","ｼｶﾞｹﾝ","ｲﾇｶﾐｸﾞﾝﾀｶﾞﾁｮｳ","ｲﾁｴﾝ","滋賀県","犬上郡多賀町","一円",0,0,0,0,0,0`},
			{`25443,"52203","5220317","ｼｶﾞｹﾝ","ｲﾇｶﾐｸﾞﾝﾀｶﾞﾁｮｳ","ｲﾁｴﾝ","滋賀県","犬上郡多賀町","一円",0,0,0,0,0,0`},
		},
		{
			{`26103,"606  ","6060017","ｷｮｳﾄﾌ","ｷｮｳﾄｼｻｷｮｳｸ","ｲﾜｸﾗｱｸﾞﾗﾁｮｳ(ｿﾉﾀ)","京都府","京都市左京区","岩倉上蔵町（その他）",1,0,0,0,0,0`},
			{`26103,"606  ","6060017","ｷｮｳﾄﾌ","ｷｮｳﾄｼｻｷｮｳｸ","ｲﾜｸﾗｱｸﾞﾗﾁｮｳ","京都府","京都市左京区","岩倉上蔵町",1,0,0,0,0,0`},
		},
		{
			{`27119,"545  ","5456090","ｵｵｻｶﾌ","ｵｵｻｶｼｱﾍﾞﾉｸ","ｱﾍﾞﾉｽｼﾞｱﾍﾞﾉﾊﾙｶｽ(ﾁｶｲ･ｶｲｿｳﾌﾒｲ)","大阪府","大阪市阿倍野区","阿倍野筋あべのハルカス（地階・階層不明）",0,0,0,0,0,0`},
			{`27119,"545  ","5456090","ｵｵｻｶﾌ","ｵｵｻｶｼｱﾍﾞﾉｸ","ｱﾍﾞﾉｽｼﾞｱﾍﾞﾉﾊﾙｶｽ","大阪府","大阪市阿倍野区","阿倍野筋あべのハルカス",0,0,0,0,0,0`},
		},
		{
			{`27119,"545  ","5450052","ｵｵｻｶﾌ","ｵｵｻｶｼｱﾍﾞﾉｸ","ｱﾍﾞﾉｽｼﾞ(ﾂｷﾞﾉﾋﾞﾙｦﾉｿﾞｸ)","大阪府","大阪市阿倍野区","阿倍野筋（次のビルを除く）",0,0,1,0,0,0`},
			{`27119,"545  ","5450052","ｵｵｻｶﾌ","ｵｵｻｶｼｱﾍﾞﾉｸ","ｱﾍﾞﾉｽｼﾞ","大阪府","大阪市阿倍野区","阿倍野筋",0,0,1,0,0,0`},
		},
		{
			{`27119,"545  ","5456060","ｵｵｻｶﾌ","ｵｵｻｶｼｱﾍﾞﾉｸ","ｱﾍﾞﾉｽｼﾞｱﾍﾞﾉﾊﾙｶｽ(60ｶｲ)","大阪府","大阪市阿倍野区","阿倍野筋あべのハルカス（６０階）",0,0,0,0,0,0`},
			{`27119,"545  ","5456060","ｵｵｻｶﾌ","ｵｵｻｶｼｱﾍﾞﾉｸ","ｱﾍﾞﾉｽｼﾞｱﾍﾞﾉﾊﾙｶｽ60ｶｲ","大阪府","大阪市阿倍野区","阿倍野筋あべのハルカス６０階",0,0,0,0,0,0`},
		},
		{
			{
				`27127,"530  ","5300041","ｵｵｻｶﾌ","ｵｵｻｶｼｷﾀｸ","ﾃﾝｼﾞﾝﾊﾞｼ(1-6ﾁｮｳﾒ)","大阪府","大阪市北区","天神橋（１～６丁目）",1,0,1,0,0,0`,
			},
			{
				`27127,"530  ","5300041","ｵｵｻｶﾌ","ｵｵｻｶｼｷﾀｸ","ﾃﾝｼﾞﾝﾊﾞｼ1ﾁｮｳﾒ","大阪府","大阪市北区","天神橋１丁目",1,0,1,0,0,0`,
				`27127,"530  ","5300041","ｵｵｻｶﾌ","ｵｵｻｶｼｷﾀｸ","ﾃﾝｼﾞﾝﾊﾞｼ2ﾁｮｳﾒ","大阪府","大阪市北区","天神橋２丁目",1,0,1,0,0,0`,
				`27127,"530  ","5300041","ｵｵｻｶﾌ","ｵｵｻｶｼｷﾀｸ","ﾃﾝｼﾞﾝﾊﾞｼ3ﾁｮｳﾒ","大阪府","大阪市北区","天神橋３丁目",1,0,1,0,0,0`,
				`27127,"530  ","5300041","ｵｵｻｶﾌ","ｵｵｻｶｼｷﾀｸ","ﾃﾝｼﾞﾝﾊﾞｼ4ﾁｮｳﾒ","大阪府","大阪市北区","天神橋４丁目",1,0,1,0,0,0`,
				`27127,"530  ","5300041","ｵｵｻｶﾌ","ｵｵｻｶｼｷﾀｸ","ﾃﾝｼﾞﾝﾊﾞｼ5ﾁｮｳﾒ","大阪府","大阪市北区","天神橋５丁目",1,0,1,0,0,0`,
				`27127,"530  ","5300041","ｵｵｻｶﾌ","ｵｵｻｶｼｷﾀｸ","ﾃﾝｼﾞﾝﾊﾞｼ6ﾁｮｳﾒ","大阪府","大阪市北区","天神橋６丁目",1,0,1,0,0,0`,
			},
		},
		{
			{
				`27127,"531  ","5310041","ｵｵｻｶﾌ","ｵｵｻｶｼｷﾀｸ","ﾃﾝｼﾞﾝﾊﾞｼ(7､8ﾁｮｳﾒ)","大阪府","大阪市北区","天神橋（７、８丁目）",1,0,1,0,0,0`,
			},
			{
				`27127,"531  ","5310041","ｵｵｻｶﾌ","ｵｵｻｶｼｷﾀｸ","ﾃﾝｼﾞﾝﾊﾞｼ7ﾁｮｳﾒ","大阪府","大阪市北区","天神橋７丁目",1,0,1,0,0,0`,
				`27127,"531  ","5310041","ｵｵｻｶﾌ","ｵｵｻｶｼｷﾀｸ","ﾃﾝｼﾞﾝﾊﾞｼ8ﾁｮｳﾒ","大阪府","大阪市北区","天神橋８丁目",1,0,1,0,0,0`,
			},
		},
		{
			{
				`01575,"04957","0495731","ﾎｯｶｲﾄﾞｳ","ｳｽｸﾞﾝｿｳﾍﾞﾂﾁｮｳ","ﾄｳﾔｺｵﾝｾﾝ(1-7ﾊﾞﾝﾁ)","北海道","有珠郡壮瞥町","洞爺湖温泉（１～７番地）",1,0,0,0,0,0`,
			},
			{
				`01575,"04957","0495731","ﾎｯｶｲﾄﾞｳ","ｳｽｸﾞﾝｿｳﾍﾞﾂﾁｮｳ","ﾄｳﾔｺｵﾝｾﾝ1ﾊﾞﾝﾁ","北海道","有珠郡壮瞥町","洞爺湖温泉１番地",1,0,0,0,0,0`,
				`01575,"04957","0495731","ﾎｯｶｲﾄﾞｳ","ｳｽｸﾞﾝｿｳﾍﾞﾂﾁｮｳ","ﾄｳﾔｺｵﾝｾﾝ2ﾊﾞﾝﾁ","北海道","有珠郡壮瞥町","洞爺湖温泉２番地",1,0,0,0,0,0`,
				`01575,"04957","0495731","ﾎｯｶｲﾄﾞｳ","ｳｽｸﾞﾝｿｳﾍﾞﾂﾁｮｳ","ﾄｳﾔｺｵﾝｾﾝ3ﾊﾞﾝﾁ","北海道","有珠郡壮瞥町","洞爺湖温泉３番地",1,0,0,0,0,0`,
				`01575,"04957","0495731","ﾎｯｶｲﾄﾞｳ","ｳｽｸﾞﾝｿｳﾍﾞﾂﾁｮｳ","ﾄｳﾔｺｵﾝｾﾝ4ﾊﾞﾝﾁ","北海道","有珠郡壮瞥町","洞爺湖温泉４番地",1,0,0,0,0,0`,
				`01575,"04957","0495731","ﾎｯｶｲﾄﾞｳ","ｳｽｸﾞﾝｿｳﾍﾞﾂﾁｮｳ","ﾄｳﾔｺｵﾝｾﾝ5ﾊﾞﾝﾁ","北海道","有珠郡壮瞥町","洞爺湖温泉５番地",1,0,0,0,0,0`,
				`01575,"04957","0495731","ﾎｯｶｲﾄﾞｳ","ｳｽｸﾞﾝｿｳﾍﾞﾂﾁｮｳ","ﾄｳﾔｺｵﾝｾﾝ6ﾊﾞﾝﾁ","北海道","有珠郡壮瞥町","洞爺湖温泉６番地",1,0,0,0,0,0`,
				`01575,"04957","0495731","ﾎｯｶｲﾄﾞｳ","ｳｽｸﾞﾝｿｳﾍﾞﾂﾁｮｳ","ﾄｳﾔｺｵﾝｾﾝ7ﾊﾞﾝﾁ","北海道","有珠郡壮瞥町","洞爺湖温泉７番地",1,0,0,0,0,0`,
			},
		},
		{
			{
				`01604,"05922","0592253","ﾎｯｶｲﾄﾞｳ","ﾆｲｶｯﾌﾟｸﾞﾝﾆｲｶｯﾌﾟﾁｮｳ","ｵｵｶﾘﾍﾞ(436､516､567ﾊﾞﾝﾁ)","北海道","新冠郡新冠町","大狩部（４３６、５１６、５６７番地）",1,0,0,0,0,0`,
			},
			{
				`01604,"05922","0592253","ﾎｯｶｲﾄﾞｳ","ﾆｲｶｯﾌﾟｸﾞﾝﾆｲｶｯﾌﾟﾁｮｳ","ｵｵｶﾘﾍﾞ436ﾊﾞﾝﾁ","北海道","新冠郡新冠町","大狩部４３６番地",1,0,0,0,0,0`,
				`01604,"05922","0592253","ﾎｯｶｲﾄﾞｳ","ﾆｲｶｯﾌﾟｸﾞﾝﾆｲｶｯﾌﾟﾁｮｳ","ｵｵｶﾘﾍﾞ516ﾊﾞﾝﾁ","北海道","新冠郡新冠町","大狩部５１６番地",1,0,0,0,0,0`,
				`01604,"05922","0592253","ﾎｯｶｲﾄﾞｳ","ﾆｲｶｯﾌﾟｸﾞﾝﾆｲｶｯﾌﾟﾁｮｳ","ｵｵｶﾘﾍﾞ567ﾊﾞﾝﾁ","北海道","新冠郡新冠町","大狩部５６７番地",1,0,0,0,0,0`,
			},
		},
		{
			{
				`44201,"870  ","8700923","ｵｵｲﾀｹﾝ","ｵｵｲﾀｼ","ﾀｶｼﾞｮｳﾆｼﾏﾁ(1-7ﾊﾞﾝ)","大分県","大分市","高城西町（１～７番）",1,0,0,0,0,0`,
			},
			{
				`44201,"870  ","8700923","ｵｵｲﾀｹﾝ","ｵｵｲﾀｼ","ﾀｶｼﾞｮｳﾆｼﾏﾁ1ﾊﾞﾝ","大分県","大分市","高城西町１番",1,0,0,0,0,0`,
				`44201,"870  ","8700923","ｵｵｲﾀｹﾝ","ｵｵｲﾀｼ","ﾀｶｼﾞｮｳﾆｼﾏﾁ2ﾊﾞﾝ","大分県","大分市","高城西町２番",1,0,0,0,0,0`,
				`44201,"870  ","8700923","ｵｵｲﾀｹﾝ","ｵｵｲﾀｼ","ﾀｶｼﾞｮｳﾆｼﾏﾁ3ﾊﾞﾝ","大分県","大分市","高城西町３番",1,0,0,0,0,0`,
				`44201,"870  ","8700923","ｵｵｲﾀｹﾝ","ｵｵｲﾀｼ","ﾀｶｼﾞｮｳﾆｼﾏﾁ4ﾊﾞﾝ","大分県","大分市","高城西町４番",1,0,0,0,0,0`,
				`44201,"870  ","8700923","ｵｵｲﾀｹﾝ","ｵｵｲﾀｼ","ﾀｶｼﾞｮｳﾆｼﾏﾁ5ﾊﾞﾝ","大分県","大分市","高城西町５番",1,0,0,0,0,0`,
				`44201,"870  ","8700923","ｵｵｲﾀｹﾝ","ｵｵｲﾀｼ","ﾀｶｼﾞｮｳﾆｼﾏﾁ6ﾊﾞﾝ","大分県","大分市","高城西町６番",1,0,0,0,0,0`,
				`44201,"870  ","8700923","ｵｵｲﾀｹﾝ","ｵｵｲﾀｼ","ﾀｶｼﾞｮｳﾆｼﾏﾁ7ﾊﾞﾝ","大分県","大分市","高城西町７番",1,0,0,0,0,0`,
			},
		},
		{
			{
				`01214,"09845","0984581","ﾎｯｶｲﾄﾞｳ","ﾜｯｶﾅｲｼ","ﾊﾞｯｶｲﾑﾗ(ｶﾐﾕｳﾁ､ｼﾓﾕｳﾁ､ﾕｳｸﾙ､ｵﾈﾄﾏﾅｲ)","北海道","稚内市","抜海村（上勇知、下勇知、夕来、オネトマナイ）",1,0,0,0,0,0`,
			},
			{
				`01214,"09845","0984581","ﾎｯｶｲﾄﾞｳ","ﾜｯｶﾅｲｼ","ﾊﾞｯｶｲﾑﾗｶﾐﾕｳﾁ","北海道","稚内市","抜海村上勇知",1,0,0,0,0,0`,
				`01214,"09845","0984581","ﾎｯｶｲﾄﾞｳ","ﾜｯｶﾅｲｼ","ﾊﾞｯｶｲﾑﾗｼﾓﾕｳﾁ","北海道","稚内市","抜海村下勇知",1,0,0,0,0,0`,
				`01214,"09845","0984581","ﾎｯｶｲﾄﾞｳ","ﾜｯｶﾅｲｼ","ﾊﾞｯｶｲﾑﾗﾕｳｸﾙ","北海道","稚内市","抜海村夕来",1,0,0,0,0,0`,
				`01214,"09845","0984581","ﾎｯｶｲﾄﾞｳ","ﾜｯｶﾅｲｼ","ﾊﾞｯｶｲﾑﾗｵﾈﾄﾏﾅｲ","北海道","稚内市","抜海村オネトマナイ",1,0,0,0,0,0`,
			},
		},
		{
			{
				`01104,"003  ","0030022","ﾎｯｶｲﾄﾞｳ","ｻｯﾎﾟﾛｼｼﾛｲｼｸ","ﾅﾝｺﾞｳﾄﾞｵﾘ(ﾐﾅﾐ)","北海道","札幌市白石区","南郷通（南）",1,0,0,0,0,0`,
			},
			{
				`01104,"003  ","0030022","ﾎｯｶｲﾄﾞｳ","ｻｯﾎﾟﾛｼｼﾛｲｼｸ","ﾅﾝｺﾞｳﾄﾞｵﾘﾐﾅﾐ","北海道","札幌市白石区","南郷通南",1,0,0,0,0,0`,
			},
		},
		{
			{
				`04101,"980  ","9800065","ﾐﾔｷﾞｹﾝ","ｾﾝﾀﾞｲｼｱｵﾊﾞｸ","ﾂﾁﾄｲ(1ﾁｮｳﾒ<11ｦﾉｿﾞｸ>)","宮城県","仙台市青葉区","土樋（１丁目「１１を除く」）",0,0,1,0,0,0`,
			},
			{
				`04101,"980  ","9800065","ﾐﾔｷﾞｹﾝ","ｾﾝﾀﾞｲｼｱｵﾊﾞｸ","ﾂﾁﾄｲ","宮城県","仙台市青葉区","土樋",0,0,1,0,0,0`,
			},
		},
		{
			{
				`01224,"066  ","0660005","ﾎｯｶｲﾄﾞｳ","ﾁﾄｾｼ","ｷｮｳﾜ(88-2､271-10､343-2､404-1､427-","北海道","千歳市","協和（８８－２、２７１－１０、３４３－２、４０４－１、４２７－",1,0,0,0,0,0`,
				`01224,"066  ","0660005","ﾎｯｶｲﾄﾞｳ","ﾁﾄｾｼ","3､431-12､443-6､608-2､641-8､814､842-","北海道","千歳市","３、４３１－１２、４４３－６、６０８－２、６４１－８、８１４、８４２－",1,0,0,0,0,0`,
				`01224,"066  ","0660005","ﾎｯｶｲﾄﾞｳ","ﾁﾄｾｼ","5､1137-3､1392､1657､1752ﾊﾞﾝﾁ)","北海道","千歳市","５、１１３７－３、１３９２、１６５７、１７５２番地）",1,0,0,0,0,0`,
				`01224,"06911","0691182","ﾎｯｶｲﾄﾞｳ","ﾁﾄｾｼ","ｷｮｳﾜ(ｿﾉﾀ)","北海道","千歳市","協和（その他）",1,0,0,0,0,0`,
			},
			{
				`01224,"066  ","0660005","ﾎｯｶｲﾄﾞｳ","ﾁﾄｾｼ","ｷｮｳﾜ","北海道","千歳市","協和",1,0,0,0,0,0`,
				`01224,"06911","0691182","ﾎｯｶｲﾄﾞｳ","ﾁﾄｾｼ","ｷｮｳﾜ","北海道","千歳市","協和",1,0,0,0,0,0`,
			},
		},
		{
			{
				`40206,"826  ","8260043","ﾌｸｵｶｹﾝ","ﾀｶﾞﾜｼ","ﾅﾗ(ｱｵﾊﾞﾁｮｳ､ｵｵｳﾗ､ｶｲｼｬﾏﾁ､ｶｽﾐｶﾞｵｶ､ｺﾞﾄｳｼﾞﾆｼﾀﾞﾝﾁ､ｺﾞﾄｳｼﾞﾋｶﾞｼﾀﾞﾝﾁ､ﾉｿﾞﾐｶﾞｵｶ､","福岡県","田川市","奈良（青葉町、大浦、会社町、霞ケ丘、後藤寺西団地、後藤寺東団地、希望ケ丘、",0,0,0,0,0,0`,
				`40206,"826  ","8260043","ﾌｸｵｶｹﾝ","ﾀｶﾞﾜｼ","ﾏﾂﾉｷ､ﾐﾂｲｺﾞﾄｳｼﾞ､ﾐﾄﾞﾘﾏﾁ､ﾂｷﾐｶﾞｵｶ)","福岡県","田川市","松の木、三井後藤寺、緑町、月見ケ丘）",0,0,0,0,0,0`,
				`40206,"826  ","8260024","ﾌｸｵｶｹﾝ","ﾀｶﾞﾜｼ","ﾆｼﾎﾝﾏﾁ","福岡県","田川市","西本町",0,0,0,0,0,0`,
			},
			{
				`40206,"826  ","8260043","ﾌｸｵｶｹﾝ","ﾀｶﾞﾜｼ","ﾅﾗｱｵﾊﾞﾁｮｳ","福岡県","田川市","奈良青葉町",0,0,0,0,0,0`,
				`40206,"826  ","8260043","ﾌｸｵｶｹﾝ","ﾀｶﾞﾜｼ","ﾅﾗｵｵｳﾗ","福岡県","田川市","奈良大浦",0,0,0,0,0,0`,
				`40206,"826  ","8260043","ﾌｸｵｶｹﾝ","ﾀｶﾞﾜｼ","ﾅﾗｶｲｼｬﾏﾁ","福岡県","田川市","奈良会社町",0,0,0,0,0,0`,
				`40206,"826  ","8260043","ﾌｸｵｶｹﾝ","ﾀｶﾞﾜｼ","ﾅﾗｶｽﾐｶﾞｵｶ","福岡県","田川市","奈良霞ケ丘",0,0,0,0,0,0`,
				`40206,"826  ","8260043","ﾌｸｵｶｹﾝ","ﾀｶﾞﾜｼ","ﾅﾗｺﾞﾄｳｼﾞﾆｼﾀﾞﾝﾁ","福岡県","田川市","奈良後藤寺西団地",0,0,0,0,0,0`,
				`40206,"826  ","8260043","ﾌｸｵｶｹﾝ","ﾀｶﾞﾜｼ","ﾅﾗｺﾞﾄｳｼﾞﾋｶﾞｼﾀﾞﾝﾁ","福岡県","田川市","奈良後藤寺東団地",0,0,0,0,0,0`,
				`40206,"826  ","8260043","ﾌｸｵｶｹﾝ","ﾀｶﾞﾜｼ","ﾅﾗﾉｿﾞﾐｶﾞｵｶ","福岡県","田川市","奈良希望ケ丘",0,0,0,0,0,0`,
				`40206,"826  ","8260043","ﾌｸｵｶｹﾝ","ﾀｶﾞﾜｼ","ﾅﾗﾏﾂﾉｷ","福岡県","田川市","奈良松の木",0,0,0,0,0,0`,
				`40206,"826  ","8260043","ﾌｸｵｶｹﾝ","ﾀｶﾞﾜｼ","ﾅﾗﾐﾂｲｺﾞﾄｳｼﾞ","福岡県","田川市","奈良三井後藤寺",0,0,0,0,0,0`,
				`40206,"826  ","8260043","ﾌｸｵｶｹﾝ","ﾀｶﾞﾜｼ","ﾅﾗﾐﾄﾞﾘﾏﾁ","福岡県","田川市","奈良緑町",0,0,0,0,0,0`,
				`40206,"826  ","8260043","ﾌｸｵｶｹﾝ","ﾀｶﾞﾜｼ","ﾅﾗﾂｷﾐｶﾞｵｶ","福岡県","田川市","奈良月見ケ丘",0,0,0,0,0,0`,
				`40206,"826  ","8260024","ﾌｸｵｶｹﾝ","ﾀｶﾞﾜｼ","ﾆｼﾎﾝﾏﾁ","福岡県","田川市","西本町",0,0,0,0,0,0`,
			},
		},
	}

	type testT struct {
		name   string
		pushes []*JapanZipCode
		pops   []*JapanZipCode
	}
	tests := []*testT{}
	for _, ba := range beforeAfter {
		tt := &testT{"", []*JapanZipCode{}, []*JapanZipCode{}}
		for _, before := range ba[0] {
			jzc, _ := parseCSV(before, false)
			tt.pushes = append(tt.pushes, jzc)
		}
		for _, after := range ba[1] {
			jzc, _ := parseCSV(after, false)
			tt.pops = append(tt.pops, jzc)
		}
		tests = append(tests, tt)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normer := newNormalizer()
			for _, push := range tt.pushes {
				normer.push(push)
			}
			for _, pop := range tt.pops {
				if got := normer.pop(); !reflect.DeepEqual(got, pop) {
					t.Errorf("normalizer.pop()  = %v, want %v", got, pop)
				}
			}
		})
	}
}

func Test_zenkaku2Int(t *testing.T) {
	type args struct {
		t string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"", args{"１２３４５６７８９０"}, 1234567890},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := zenkaku2Int(tt.args.t); got != tt.want {
				t.Errorf("zenkaku2Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_int2Zenkaku(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{1234567890}, "１２３４５６７８９０"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := int2Zenkaku(tt.args.i); got != tt.want {
				t.Errorf("int2Zenkaku() = %v, want %v", got, tt.want)
			}
		})
	}
}
