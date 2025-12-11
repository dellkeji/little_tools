// 扩展的学习内容数据库

const learningData = {
    english: {
        lessons: [
            // 水果类 (Fruits)
            { word: 'Apple', translation: '苹果', emoji: '🍎', category: '水果' },
            { word: 'Banana', translation: '香蕉', emoji: '🍌', category: '水果' },
            { word: 'Orange', translation: '橙子', emoji: '🍊', category: '水果' },
            { word: 'Grape', translation: '葡萄', emoji: '🍇', category: '水果' },
            { word: 'Watermelon', translation: '西瓜', emoji: '🍉', category: '水果' },
            { word: 'Strawberry', translation: '草莓', emoji: '🍓', category: '水果' },
            { word: 'Peach', translation: '桃子', emoji: '🍑', category: '水果' },
            { word: 'Pear', translation: '梨', emoji: '🍐', category: '水果' },
            
            // 动物类 (Animals)
            { word: 'Cat', translation: '猫', emoji: '🐱', category: '动物' },
            { word: 'Dog', translation: '狗', emoji: '🐶', category: '动物' },
            { word: 'Bird', translation: '鸟', emoji: '🐦', category: '动物' },
            { word: 'Fish', translation: '鱼', emoji: '🐟', category: '动物' },
            { word: 'Elephant', translation: '大象', emoji: '🐘', category: '动物' },
            { word: 'Lion', translation: '狮子', emoji: '🦁', category: '动物' },
            { word: 'Tiger', translation: '老虎', emoji: '🐯', category: '动物' },
            { word: 'Rabbit', translation: '兔子', emoji: '🐰', category: '动物' },
            { word: 'Monkey', translation: '猴子', emoji: '🐵', category: '动物' },
            { word: 'Panda', translation: '熊猫', emoji: '🐼', category: '动物' },
            
            // 颜色类 (Colors)
            { word: 'Red', translation: '红色', emoji: '🔴', category: '颜色' },
            { word: 'Blue', translation: '蓝色', emoji: '🔵', category: '颜色' },
            { word: 'Yellow', translation: '黄色', emoji: '🟡', category: '颜色' },
            { word: 'Green', translation: '绿色', emoji: '🟢', category: '颜色' },
            { word: 'Orange', translation: '橙色', emoji: '🟠', category: '颜色' },
            { word: 'Purple', translation: '紫色', emoji: '🟣', category: '颜色' },
            { word: 'Black', translation: '黑色', emoji: '⚫', category: '颜色' },
            { word: 'White', translation: '白色', emoji: '⚪', category: '颜色' },
            
            // 数字类 (Numbers)
            { word: 'One', translation: '一', emoji: '1️⃣', category: '数字' },
            { word: 'Two', translation: '二', emoji: '2️⃣', category: '数字' },
            { word: 'Three', translation: '三', emoji: '3️⃣', category: '数字' },
            { word: 'Four', translation: '四', emoji: '4️⃣', category: '数字' },
            { word: 'Five', translation: '五', emoji: '5️⃣', category: '数字' },
            { word: 'Six', translation: '六', emoji: '6️⃣', category: '数字' },
            { word: 'Seven', translation: '七', emoji: '7️⃣', category: '数字' },
            { word: 'Eight', translation: '八', emoji: '8️⃣', category: '数字' },
            { word: 'Nine', translation: '九', emoji: '9️⃣', category: '数字' },
            { word: 'Ten', translation: '十', emoji: '🔟', category: '数字' },
            
            // 身体部位 (Body Parts)
            { word: 'Eye', translation: '眼睛', emoji: '👁️', category: '身体' },
            { word: 'Ear', translation: '耳朵', emoji: '👂', category: '身体' },
            { word: 'Nose', translation: '鼻子', emoji: '👃', category: '身体' },
            { word: 'Mouth', translation: '嘴巴', emoji: '👄', category: '身体' },
            { word: 'Hand', translation: '手', emoji: '✋', category: '身体' },
            { word: 'Foot', translation: '脚', emoji: '🦶', category: '身体' },
            
            // 日常用品 (Daily Items)
            { word: 'Book', translation: '书', emoji: '📖', category: '用品' },
            { word: 'Pen', translation: '笔', emoji: '✒️', category: '用品' },
            { word: 'Ball', translation: '球', emoji: '⚽', category: '用品' },
            { word: 'Car', translation: '汽车', emoji: '🚗', category: '用品' },
            { word: 'House', translation: '房子', emoji: '🏠', category: '用品' },
            { word: 'Tree', translation: '树', emoji: '🌳', category: '用品' },
            { word: 'Sun', translation: '太阳', emoji: '☀️', category: '用品' },
            { word: 'Moon', translation: '月亮', emoji: '🌙', category: '用品' },
            { word: 'Star', translation: '星星', emoji: '⭐', category: '用品' }
        ],
        tests: [
            // 水果测试
            { question: 'Apple 的中文意思是？', options: ['苹果', '香蕉', '橙子', '梨'], answer: 0 },
            { question: 'Banana 的中文意思是？', options: ['苹果', '香蕉', '橙子', '葡萄'], answer: 1 },
            { question: '🍊 对应的英文是？', options: ['Apple', 'Banana', 'Orange', 'Grape'], answer: 2 },
            { question: '草莓 的英文是？', options: ['Peach', 'Pear', 'Strawberry', 'Watermelon'], answer: 2 },
            
            // 动物测试
            { question: 'Cat 的中文意思是？', options: ['狗', '猫', '鸟', '鱼'], answer: 1 },
            { question: 'Dog 的中文意思是？', options: ['猫', '狗', '兔子', '猴子'], answer: 1 },
            { question: '🐘 对应的英文是？', options: ['Lion', 'Tiger', 'Elephant', 'Panda'], answer: 2 },
            { question: '熊猫 的英文是？', options: ['Monkey', 'Rabbit', 'Tiger', 'Panda'], answer: 3 },
            { question: 'Lion 的中文意思是？', options: ['老虎', '狮子', '猴子', '兔子'], answer: 1 },
            
            // 颜色测试
            { question: 'Red 的中文意思是？', options: ['红色', '蓝色', '黄色', '绿色'], answer: 0 },
            { question: 'Blue 的中文意思是？', options: ['红色', '蓝色', '黄色', '绿色'], answer: 1 },
            { question: '🟡 对应的英文是？', options: ['Red', 'Blue', 'Yellow', 'Green'], answer: 2 },
            { question: '绿色 的英文是？', options: ['Red', 'Blue', 'Yellow', 'Green'], answer: 3 },
            
            // 数字测试
            { question: 'One 的中文意思是？', options: ['一', '二', '三', '四'], answer: 0 },
            { question: 'Five 的中文意思是？', options: ['三', '四', '五', '六'], answer: 2 },
            { question: '十 的英文是？', options: ['Eight', 'Nine', 'Ten', 'Seven'], answer: 2 },
            
            // 身体部位测试
            { question: 'Eye 的中文意思是？', options: ['眼睛', '耳朵', '鼻子', '嘴巴'], answer: 0 },
            { question: 'Hand 的中文意思是？', options: ['脚', '手', '眼睛', '耳朵'], answer: 1 },
            
            // 日常用品测试
            { question: 'Book 的中文意思是？', options: ['书', '笔', '球', '车'], answer: 0 },
            { question: '⚽ 对应的英文是？', options: ['Book', 'Pen', 'Ball', 'Car'], answer: 2 }
        ]
    },
    math: {
        lessons: [
            // 数字认知
            { title: '认识数字 1-5', content: '1️⃣ 一  2️⃣ 二  3️⃣ 三  4️⃣ 四  5️⃣ 五', emoji: '🔢' },
            { title: '认识数字 6-10', content: '6️⃣ 六  7️⃣ 七  8️⃣ 八  9️⃣ 九  🔟 十', emoji: '🔢' },
            { title: '数字 0', content: '0️⃣ 零 - 表示什么都没有', emoji: '⭕' },
            
            // 加法
            { title: '加法：1+1', content: '🍎 + 🍎 = 🍎🍎  (1 + 1 = 2)', emoji: '➕' },
            { title: '加法：2+1', content: '🍎🍎 + 🍎 = 🍎🍎🍎  (2 + 1 = 3)', emoji: '➕' },
            { title: '加法：2+2', content: '🍎🍎 + 🍎🍎 = 🍎🍎🍎🍎  (2 + 2 = 4)', emoji: '➕' },
            { title: '加法：3+2', content: '🍎🍎🍎 + 🍎🍎 = 🍎🍎🍎🍎🍎  (3 + 2 = 5)', emoji: '➕' },
            { title: '加法：4+1', content: '🍎🍎🍎🍎 + 🍎 = 🍎🍎🍎🍎🍎  (4 + 1 = 5)', emoji: '➕' },
            { title: '加法：5+5', content: '✋ + ✋ = 🔟  (5 + 5 = 10)', emoji: '➕' },
            
            // 减法
            { title: '减法：2-1', content: '🍎🍎 - 🍎 = 🍎  (2 - 1 = 1)', emoji: '➖' },
            { title: '减法：3-1', content: '🍎🍎🍎 - 🍎 = 🍎🍎  (3 - 1 = 2)', emoji: '➖' },
            { title: '减法：4-2', content: '🍎🍎🍎🍎 - 🍎🍎 = 🍎🍎  (4 - 2 = 2)', emoji: '➖' },
            { title: '减法：5-3', content: '🍎🍎🍎🍎🍎 - 🍎🍎🍎 = 🍎🍎  (5 - 3 = 2)', emoji: '➖' },
            { title: '减法：10-5', content: '🔟 - ✋ = ✋  (10 - 5 = 5)', emoji: '➖' },
            
            // 形状
            { title: '圆形', content: '⭕ 圆圆的，像球一样', emoji: '⭕' },
            { title: '正方形', content: '⬜ 四条边一样长', emoji: '⬜' },
            { title: '三角形', content: '🔺 有三个角', emoji: '🔺' },
            { title: '长方形', content: '▭ 对边一样长', emoji: '▭' },
            { title: '星形', content: '⭐ 像星星一样', emoji: '⭐' },
            { title: '心形', content: '❤️ 像爱心一样', emoji: '❤️' },
            
            // 比较大小
            { title: '比大小：多与少', content: '🍎🍎🍎 > 🍎  (3个 多于 1个)', emoji: '🔍' },
            { title: '比大小：相等', content: '🍎🍎 = 🍎🍎  (2个 等于 2个)', emoji: '🔍' },
            { title: '比大小：大与小', content: '🐘 大  🐭 小', emoji: '🔍' },
            { title: '比大小：高与矮', content: '🦒 高  🐶 矮', emoji: '🔍' }
        ],
        tests: [
            // 加法测试
            { question: '1 + 1 = ?', options: ['1', '2', '3', '4'], answer: 1 },
            { question: '2 + 1 = ?', options: ['2', '3', '4', '5'], answer: 1 },
            { question: '2 + 2 = ?', options: ['2', '3', '4', '5'], answer: 2 },
            { question: '3 + 1 = ?', options: ['3', '4', '5', '6'], answer: 1 },
            { question: '3 + 2 = ?', options: ['4', '5', '6', '7'], answer: 1 },
            { question: '4 + 1 = ?', options: ['4', '5', '6', '7'], answer: 1 },
            { question: '5 + 5 = ?', options: ['8', '9', '10', '11'], answer: 2 },
            { question: '1 + 2 = ?', options: ['2', '3', '4', '5'], answer: 1 },
            { question: '4 + 2 = ?', options: ['5', '6', '7', '8'], answer: 1 },
            { question: '3 + 3 = ?', options: ['5', '6', '7', '8'], answer: 1 },
            
            // 减法测试
            { question: '2 - 1 = ?', options: ['0', '1', '2', '3'], answer: 1 },
            { question: '3 - 1 = ?', options: ['1', '2', '3', '4'], answer: 1 },
            { question: '4 - 2 = ?', options: ['1', '2', '3', '4'], answer: 1 },
            { question: '5 - 2 = ?', options: ['2', '3', '4', '5'], answer: 1 },
            { question: '5 - 3 = ?', options: ['1', '2', '3', '4'], answer: 1 },
            { question: '10 - 5 = ?', options: ['3', '4', '5', '6'], answer: 2 },
            { question: '6 - 3 = ?', options: ['2', '3', '4', '5'], answer: 1 },
            { question: '7 - 2 = ?', options: ['4', '5', '6', '7'], answer: 1 },
            
            // 形状测试
            { question: '⭕ 是什么形状？', options: ['圆形', '正方形', '三角形', '长方形'], answer: 0 },
            { question: '⬜ 是什么形状？', options: ['圆形', '正方形', '三角形', '长方形'], answer: 1 },
            { question: '🔺 是什么形状？', options: ['圆形', '正方形', '三角形', '长方形'], answer: 2 },
            { question: '⭐ 是什么形状？', options: ['圆形', '星形', '三角形', '心形'], answer: 1 },
            
            // 比较大小测试
            { question: '3 和 1 哪个大？', options: ['1', '3', '一样大', '不知道'], answer: 1 },
            { question: '5 和 8 哪个小？', options: ['5', '8', '一样大', '不知道'], answer: 0 },
            { question: '2 + 2 和 4 比较', options: ['2+2大', '4大', '一样大', '不知道'], answer: 2 },
            { question: '🍎🍎🍎 和 🍎🍎 哪个多？', options: ['左边多', '右边多', '一样多', '不知道'], answer: 0 }
        ]
    },
    chinese: {
        lessons: [
            // 数字汉字
            { char: '一', pinyin: 'yī', meaning: '数字1', category: '数字' },
            { char: '二', pinyin: 'èr', meaning: '数字2', category: '数字' },
            { char: '三', pinyin: 'sān', meaning: '数字3', category: '数字' },
            { char: '四', pinyin: 'sì', meaning: '数字4', category: '数字' },
            { char: '五', pinyin: 'wǔ', meaning: '数字5', category: '数字' },
            { char: '六', pinyin: 'liù', meaning: '数字6', category: '数字' },
            { char: '七', pinyin: 'qī', meaning: '数字7', category: '数字' },
            { char: '八', pinyin: 'bā', meaning: '数字8', category: '数字' },
            { char: '九', pinyin: 'jiǔ', meaning: '数字9', category: '数字' },
            { char: '十', pinyin: 'shí', meaning: '数字10', category: '数字' },
            
            // 方位词
            { char: '上', pinyin: 'shàng', meaning: '上面', category: '方位' },
            { char: '下', pinyin: 'xià', meaning: '下面', category: '方位' },
            { char: '左', pinyin: 'zuǒ', meaning: '左边', category: '方位' },
            { char: '右', pinyin: 'yòu', meaning: '右边', category: '方位' },
            { char: '前', pinyin: 'qián', meaning: '前面', category: '方位' },
            { char: '后', pinyin: 'hòu', meaning: '后面', category: '方位' },
            { char: '中', pinyin: 'zhōng', meaning: '中间', category: '方位' },
            
            // 形容词
            { char: '大', pinyin: 'dà', meaning: '大的', category: '形容词' },
            { char: '小', pinyin: 'xiǎo', meaning: '小的', category: '形容词' },
            { char: '多', pinyin: 'duō', meaning: '很多', category: '形容词' },
            { char: '少', pinyin: 'shǎo', meaning: '很少', category: '形容词' },
            { char: '长', pinyin: 'cháng', meaning: '长的', category: '形容词' },
            { char: '短', pinyin: 'duǎn', meaning: '短的', category: '形容词' },
            { char: '高', pinyin: 'gāo', meaning: '高的', category: '形容词' },
            { char: '矮', pinyin: 'ǎi', meaning: '矮的', category: '形容词' },
            { char: '好', pinyin: 'hǎo', meaning: '好的', category: '形容词' },
            { char: '坏', pinyin: 'huài', meaning: '坏的', category: '形容词' },
            
            // 家庭成员
            { char: '爸', pinyin: 'bà', meaning: '爸爸', category: '家庭' },
            { char: '妈', pinyin: 'mā', meaning: '妈妈', category: '家庭' },
            { char: '哥', pinyin: 'gē', meaning: '哥哥', category: '家庭' },
            { char: '姐', pinyin: 'jiě', meaning: '姐姐', category: '家庭' },
            { char: '弟', pinyin: 'dì', meaning: '弟弟', category: '家庭' },
            { char: '妹', pinyin: 'mèi', meaning: '妹妹', category: '家庭' },
            
            // 常用字
            { char: '人', pinyin: 'rén', meaning: '人类', category: '常用' },
            { char: '口', pinyin: 'kǒu', meaning: '嘴巴', category: '常用' },
            { char: '手', pinyin: 'shǒu', meaning: '手', category: '常用' },
            { char: '足', pinyin: 'zú', meaning: '脚', category: '常用' },
            { char: '目', pinyin: 'mù', meaning: '眼睛', category: '常用' },
            { char: '耳', pinyin: 'ěr', meaning: '耳朵', category: '常用' },
            { char: '日', pinyin: 'rì', meaning: '太阳/日子', category: '常用' },
            { char: '月', pinyin: 'yuè', meaning: '月亮/月份', category: '常用' },
            { char: '水', pinyin: 'shuǐ', meaning: '水', category: '常用' },
            { char: '火', pinyin: 'huǒ', meaning: '火', category: '常用' },
            { char: '山', pinyin: 'shān', meaning: '山', category: '常用' },
            { char: '石', pinyin: 'shí', meaning: '石头', category: '常用' },
            { char: '田', pinyin: 'tián', meaning: '田地', category: '常用' },
            { char: '土', pinyin: 'tǔ', meaning: '土地', category: '常用' }
        ],
        tests: [
            // 数字测试
            { question: '"一" 的拼音是？', options: ['yī', 'èr', 'sān', 'sì'], answer: 0 },
            { question: '"五" 的拼音是？', options: ['sì', 'wǔ', 'liù', 'qī'], answer: 1 },
            { question: '拼音 "sān" 对应的汉字是？', options: ['一', '二', '三', '四'], answer: 2 },
            { question: '拼音 "shí" 对应的汉字是？', options: ['七', '八', '九', '十'], answer: 3 },
            { question: '"八" 的拼音是？', options: ['liù', 'qī', 'bā', 'jiǔ'], answer: 2 },
            
            // 方位词测试
            { question: '"上" 的拼音是？', options: ['shàng', 'xià', 'zuǒ', 'yòu'], answer: 0 },
            { question: '"下" 的拼音是？', options: ['shàng', 'xià', 'zuǒ', 'yòu'], answer: 1 },
            { question: '拼音 "zuǒ" 对应的汉字是？', options: ['上', '下', '左', '右'], answer: 2 },
            { question: '"中" 的拼音是？', options: ['qián', 'hòu', 'zhōng', 'wài'], answer: 2 },
            
            // 形容词测试
            { question: '"大" 的拼音是？', options: ['xiǎo', 'dà', 'duō', 'shǎo'], answer: 1 },
            { question: '"小" 的拼音是？', options: ['dà', 'xiǎo', 'duō', 'shǎo'], answer: 1 },
            { question: '拼音 "gāo" 对应的汉字是？', options: ['大', '小', '高', '矮'], answer: 2 },
            { question: '"好" 的拼音是？', options: ['hǎo', 'huài', 'duō', 'shǎo'], answer: 0 },
            
            // 家庭成员测试
            { question: '"爸" 的拼音是？', options: ['bà', 'mā', 'gē', 'jiě'], answer: 0 },
            { question: '"妈" 的拼音是？', options: ['bà', 'mā', 'gē', 'jiě'], answer: 1 },
            { question: '拼音 "gē" 对应的汉字是？', options: ['爸', '妈', '哥', '姐'], answer: 2 },
            
            // 常用字测试
            { question: '"人" 的拼音是？', options: ['rén', 'kǒu', 'shǒu', 'zú'], answer: 0 },
            { question: '"日" 的拼音是？', options: ['rì', 'yuè', 'shuǐ', 'huǒ'], answer: 0 },
            { question: '拼音 "shuǐ" 对应的汉字是？', options: ['日', '月', '水', '火'], answer: 2 },
            { question: '"山" 的拼音是？', options: ['shān', 'shí', 'tián', 'tǔ'], answer: 0 }
        ]
    }
};


// 扩展英语单词库 - 更多分类
learningData.english.lessons.push(
    // 食物类 (Food)
    { word: 'Bread', translation: '面包', emoji: '🍞', category: '食物' },
    { word: 'Milk', translation: '牛奶', emoji: '🥛', category: '食物' },
    { word: 'Egg', translation: '鸡蛋', emoji: '🥚', category: '食物' },
    { word: 'Cake', translation: '蛋糕', emoji: '🍰', category: '食物' },
    { word: 'Pizza', translation: '披萨', emoji: '🍕', category: '食物' },
    { word: 'Rice', translation: '米饭', emoji: '🍚', category: '食物' },
    { word: 'Noodle', translation: '面条', emoji: '🍜', category: '食物' },
    { word: 'Candy', translation: '糖果', emoji: '🍬', category: '食物' },
    { word: 'Cookie', translation: '饼干', emoji: '🍪', category: '食物' },
    { word: 'Juice', translation: '果汁', emoji: '🧃', category: '食物' },
    
    // 交通工具 (Transportation)
    { word: 'Bus', translation: '公交车', emoji: '🚌', category: '交通' },
    { word: 'Train', translation: '火车', emoji: '🚂', category: '交通' },
    { word: 'Plane', translation: '飞机', emoji: '✈️', category: '交通' },
    { word: 'Bike', translation: '自行车', emoji: '🚲', category: '交通' },
    { word: 'Ship', translation: '轮船', emoji: '🚢', category: '交通' },
    { word: 'Taxi', translation: '出租车', emoji: '🚕', category: '交通' },
    
    // 天气 (Weather)
    { word: 'Rain', translation: '雨', emoji: '🌧️', category: '天气' },
    { word: 'Snow', translation: '雪', emoji: '❄️', category: '天气' },
    { word: 'Wind', translation: '风', emoji: '💨', category: '天气' },
    { word: 'Cloud', translation: '云', emoji: '☁️', category: '天气' },
    
    // 衣服 (Clothes)
    { word: 'Hat', translation: '帽子', emoji: '🎩', category: '衣服' },
    { word: 'Shirt', translation: '衬衫', emoji: '👔', category: '衣服' },
    { word: 'Shoes', translation: '鞋子', emoji: '👞', category: '衣服' },
    { word: 'Dress', translation: '连衣裙', emoji: '👗', category: '衣服' },
    { word: 'Pants', translation: '裤子', emoji: '👖', category: '衣服' },
    { word: 'Socks', translation: '袜子', emoji: '🧦', category: '衣服' },
    
    // 运动 (Sports)
    { word: 'Football', translation: '足球', emoji: '⚽', category: '运动' },
    { word: 'Basketball', translation: '篮球', emoji: '🏀', category: '运动' },
    { word: 'Swimming', translation: '游泳', emoji: '🏊', category: '运动' },
    { word: 'Running', translation: '跑步', emoji: '🏃', category: '运动' },
    
    // 乐器 (Music)
    { word: 'Piano', translation: '钢琴', emoji: '🎹', category: '乐器' },
    { word: 'Guitar', translation: '吉他', emoji: '🎸', category: '乐器' },
    { word: 'Drum', translation: '鼓', emoji: '🥁', category: '乐器' }
);

// 扩展英语测试题
learningData.english.tests.push(
    // 食物测试
    { question: 'Bread 的中文意思是？', options: ['面包', '牛奶', '鸡蛋', '蛋糕'], answer: 0 },
    { question: 'Milk 的中文意思是？', options: ['面包', '牛奶', '鸡蛋', '蛋糕'], answer: 1 },
    { question: '🥚 对应的英文是？', options: ['Bread', 'Milk', 'Egg', 'Cake'], answer: 2 },
    { question: 'Pizza 的中文意思是？', options: ['米饭', '面条', '披萨', '糖果'], answer: 2 },
    { question: '🍜 对应的英文是？', options: ['Rice', 'Noodle', 'Candy', 'Cookie'], answer: 1 },
    { question: 'Juice 的中文意思是？', options: ['牛奶', '果汁', '水', '茶'], answer: 1 },
    
    // 交通工具测试
    { question: 'Bus 的中文意思是？', options: ['公交车', '火车', '飞机', '自行车'], answer: 0 },
    { question: '🚂 对应的英文是？', options: ['Bus', 'Train', 'Plane', 'Bike'], answer: 1 },
    { question: 'Plane 的中文意思是？', options: ['公交车', '火车', '飞机', '轮船'], answer: 2 },
    { question: '🚲 对应的英文是？', options: ['Car', 'Bus', 'Bike', 'Ship'], answer: 2 },
    
    // 天气测试
    { question: 'Rain 的中文意思是？', options: ['雨', '雪', '风', '云'], answer: 0 },
    { question: '❄️ 对应的英文是？', options: ['Rain', 'Snow', 'Wind', 'Cloud'], answer: 1 },
    { question: 'Wind 的中文意思是？', options: ['雨', '雪', '风', '云'], answer: 2 },
    
    // 衣服测试
    { question: 'Hat 的中文意思是？', options: ['帽子', '衬衫', '鞋子', '裤子'], answer: 0 },
    { question: '👗 对应的英文是？', options: ['Hat', 'Shirt', 'Dress', 'Pants'], answer: 2 },
    { question: 'Shoes 的中文意思是？', options: ['帽子', '衬衫', '鞋子', '袜子'], answer: 2 },
    
    // 运动测试
    { question: 'Football 的中文意思是？', options: ['足球', '篮球', '游泳', '跑步'], answer: 0 },
    { question: '🏀 对应的英文是？', options: ['Football', 'Basketball', 'Swimming', 'Running'], answer: 1 },
    
    // 乐器测试
    { question: 'Piano 的中文意思是？', options: ['钢琴', '吉他', '鼓', '小提琴'], answer: 0 },
    { question: '🎸 对应的英文是？', options: ['Piano', 'Guitar', 'Drum', 'Violin'], answer: 1 }
);

// 扩展数学题库
learningData.math.tests.push(
    // 更多加法
    { question: '6 + 1 = ?', options: ['5', '6', '7', '8'], answer: 2 },
    { question: '6 + 2 = ?', options: ['6', '7', '8', '9'], answer: 2 },
    { question: '7 + 1 = ?', options: ['6', '7', '8', '9'], answer: 2 },
    { question: '7 + 2 = ?', options: ['7', '8', '9', '10'], answer: 2 },
    { question: '8 + 1 = ?', options: ['7', '8', '9', '10'], answer: 2 },
    { question: '8 + 2 = ?', options: ['8', '9', '10', '11'], answer: 2 },
    { question: '9 + 1 = ?', options: ['8', '9', '10', '11'], answer: 2 },
    { question: '4 + 4 = ?', options: ['6', '7', '8', '9'], answer: 2 },
    { question: '6 + 4 = ?', options: ['8', '9', '10', '11'], answer: 2 },
    { question: '7 + 3 = ?', options: ['8', '9', '10', '11'], answer: 2 },
    
    // 更多减法
    { question: '8 - 3 = ?', options: ['4', '5', '6', '7'], answer: 1 },
    { question: '9 - 4 = ?', options: ['4', '5', '6', '7'], answer: 1 },
    { question: '10 - 3 = ?', options: ['5', '6', '7', '8'], answer: 2 },
    { question: '10 - 4 = ?', options: ['5', '6', '7', '8'], answer: 1 },
    { question: '9 - 5 = ?', options: ['3', '4', '5', '6'], answer: 1 },
    { question: '8 - 4 = ?', options: ['3', '4', '5', '6'], answer: 1 },
    { question: '7 - 3 = ?', options: ['3', '4', '5', '6'], answer: 1 },
    { question: '6 - 4 = ?', options: ['1', '2', '3', '4'], answer: 1 },
    { question: '9 - 3 = ?', options: ['5', '6', '7', '8'], answer: 1 },
    { question: '10 - 6 = ?', options: ['3', '4', '5', '6'], answer: 1 }
);

// 扩展语文汉字库
learningData.chinese.lessons.push(
    // 更多常用字
    { char: '天', pinyin: 'tiān', meaning: '天空', category: '常用' },
    { char: '地', pinyin: 'dì', meaning: '地面', category: '常用' },
    { char: '木', pinyin: 'mù', meaning: '树木', category: '常用' },
    { char: '林', pinyin: 'lín', meaning: '树林', category: '常用' },
    { char: '森', pinyin: 'sēn', meaning: '森林', category: '常用' },
    { char: '雨', pinyin: 'yǔ', meaning: '下雨', category: '常用' },
    { char: '雪', pinyin: 'xuě', meaning: '雪花', category: '常用' },
    { char: '风', pinyin: 'fēng', meaning: '风', category: '常用' },
    { char: '云', pinyin: 'yún', meaning: '云朵', category: '常用' },
    { char: '电', pinyin: 'diàn', meaning: '电', category: '常用' },
    
    // 动物字
    { char: '牛', pinyin: 'niú', meaning: '牛', category: '动物' },
    { char: '羊', pinyin: 'yáng', meaning: '羊', category: '动物' },
    { char: '马', pinyin: 'mǎ', meaning: '马', category: '动物' },
    { char: '鸟', pinyin: 'niǎo', meaning: '鸟', category: '动物' },
    { char: '鱼', pinyin: 'yú', meaning: '鱼', category: '动物' },
    { char: '虫', pinyin: 'chóng', meaning: '虫子', category: '动物' },
    
    // 颜色字
    { char: '红', pinyin: 'hóng', meaning: '红色', category: '颜色' },
    { char: '黄', pinyin: 'huáng', meaning: '黄色', category: '颜色' },
    { char: '蓝', pinyin: 'lán', meaning: '蓝色', category: '颜色' },
    { char: '绿', pinyin: 'lǜ', meaning: '绿色', category: '颜色' },
    { char: '白', pinyin: 'bái', meaning: '白色', category: '颜色' },
    { char: '黑', pinyin: 'hēi', meaning: '黑色', category: '颜色' },
    
    // 动词
    { char: '来', pinyin: 'lái', meaning: '来', category: '动词' },
    { char: '去', pinyin: 'qù', meaning: '去', category: '动词' },
    { char: '看', pinyin: 'kàn', meaning: '看', category: '动词' },
    { char: '听', pinyin: 'tīng', meaning: '听', category: '动词' },
    { char: '说', pinyin: 'shuō', meaning: '说话', category: '动词' },
    { char: '吃', pinyin: 'chī', meaning: '吃', category: '动词' },
    { char: '喝', pinyin: 'hē', meaning: '喝', category: '动词' },
    { char: '走', pinyin: 'zǒu', meaning: '走路', category: '动词' },
    { char: '跑', pinyin: 'pǎo', meaning: '跑步', category: '动词' },
    { char: '跳', pinyin: 'tiào', meaning: '跳', category: '动词' }
);

// 扩展语文测试题
learningData.chinese.tests.push(
    // 常用字测试
    { question: '"天" 的拼音是？', options: ['tiān', 'dì', 'mù', 'lín'], answer: 0 },
    { question: '"地" 的拼音是？', options: ['tiān', 'dì', 'mù', 'lín'], answer: 1 },
    { question: '拼音 "yǔ" 对应的汉字是？', options: ['天', '地', '雨', '雪'], answer: 2 },
    { question: '"风" 的拼音是？', options: ['fēng', 'yún', 'diàn', 'shuǐ'], answer: 0 },
    { question: '拼音 "xuě" 对应的汉字是？', options: ['雨', '雪', '风', '云'], answer: 1 },
    
    // 动物字测试
    { question: '"牛" 的拼音是？', options: ['niú', 'yáng', 'mǎ', 'niǎo'], answer: 0 },
    { question: '"羊" 的拼音是？', options: ['niú', 'yáng', 'mǎ', 'niǎo'], answer: 1 },
    { question: '拼音 "yú" 对应的汉字是？', options: ['牛', '羊', '鱼', '鸟'], answer: 2 },
    { question: '"马" 的拼音是？', options: ['niú', 'yáng', 'mǎ', 'yú'], answer: 2 },
    
    // 颜色字测试
    { question: '"红" 的拼音是？', options: ['hóng', 'huáng', 'lán', 'lǜ'], answer: 0 },
    { question: '"黄" 的拼音是？', options: ['hóng', 'huáng', 'lán', 'lǜ'], answer: 1 },
    { question: '拼音 "lán" 对应的汉字是？', options: ['红', '黄', '蓝', '绿'], answer: 2 },
    { question: '"白" 的拼音是？', options: ['bái', 'hēi', 'hóng', 'lǜ'], answer: 0 },
    
    // 动词测试
    { question: '"来" 的拼音是？', options: ['lái', 'qù', 'kàn', 'tīng'], answer: 0 },
    { question: '"去" 的拼音是？', options: ['lái', 'qù', 'kàn', 'tīng'], answer: 1 },
    { question: '拼音 "chī" 对应的汉字是？', options: ['看', '听', '吃', '喝'], answer: 2 },
    { question: '"跑" 的拼音是？', options: ['zǒu', 'pǎo', 'tiào', 'fēi'], answer: 1 },
    { question: '拼音 "shuō" 对应的汉字是？', options: ['看', '听', '说', '走'], answer: 2 }
);
