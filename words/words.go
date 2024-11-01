package example_words

// Word represents a vocabulary entry
type Word struct {
	Pinyin     string
	Chinese    string
	Definition string
	Example    string
}

var ExampleWords = []Word{
	{
		Pinyin:     "nǐ hǎo",
		Chinese:    "你好",
		Definition: "Hello",
		Example:    "",
	},
	{
		Pinyin:     "xiè xiè",
		Chinese:    "谢谢",
		Definition: "Thank you",
		Example:    "",
	},
	{
		Pinyin:     "zǎo",
		Chinese:    "早",
		Definition: "Morning",
		Example:    "Zǎo, good morning!",
	},
	{
		Pinyin:     "péngyou",
		Chinese:    "朋友",
		Definition: "Friend",
		Example:    "Wǒ de péngyou hěn hǎo.",
	},
	{
		Pinyin:     "chī fàn",
		Chinese:    "吃饭",
		Definition: "Eat meal",
		Example:    "Wǒmen qù chī fàn.",
	},
	{
		Pinyin:     "hǎo",
		Chinese:    "好",
		Definition: "Good",
		Example:    "Hěn hǎo, that's good!",
	},
	{
		Pinyin:     "shuǐ",
		Chinese:    "水",
		Definition: "Water",
		Example:    "Wǒ yào yī bēi shuǐ.",
	},
	{
		Pinyin:     "ài",
		Chinese:    "爱",
		Definition: "Love",
		Example:    "Wǒ ài nǐ means I love you.",
	},
	{
		Pinyin:     "rén",
		Chinese:    "人",
		Definition: "Person",
		Example:    "Měi gè rén dōu bù tóng.",
	},
	{
		Pinyin:     "jiā",
		Chinese:    "家",
		Definition: "Home/Family",
		Example:    "Wǒ de jiā zài Běijīng.",
	},
}
