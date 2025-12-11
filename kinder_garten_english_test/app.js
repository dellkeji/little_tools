// 学习数据从 learningData.js 加载

let currentTest = null;
let currentQuestionIndex = 0;
let score = 0;
let currentLessons = []; // 当前显示的学习内容
let currentTests = []; // 当前测试题目

// 工具函数：随机打乱数组
function shuffleArray(array) {
    const newArray = [...array];
    for (let i = newArray.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [newArray[i], newArray[j]] = [newArray[j], newArray[i]];
    }
    return newArray;
}

// 工具函数：从数组中随机选择n个元素
function getRandomItems(array, count) {
    const shuffled = shuffleArray(array);
    return shuffled.slice(0, Math.min(count, array.length));
}

// 页面导航
function showHome() {
    document.querySelectorAll('.page').forEach(page => page.classList.remove('active'));
    document.getElementById('home-page').classList.add('active');
}

function showSubject(subject) {
    document.querySelectorAll('.page').forEach(page => page.classList.remove('active'));
    document.getElementById(subject + '-page').classList.add('active');
    document.getElementById(subject + '-content').innerHTML = '';
}

// 语音合成功能
function speak(text, lang = 'en-US') {
    // 检查浏览器是否支持语音合成
    if ('speechSynthesis' in window) {
        // 停止当前正在播放的语音
        window.speechSynthesis.cancel();
        
        const utterance = new SpeechSynthesisUtterance(text);
        utterance.lang = lang;
        utterance.rate = 0.8; // 语速稍慢，适合儿童学习
        utterance.pitch = 1.1; // 音调稍高，更活泼
        utterance.volume = 1.0; // 音量
        
        window.speechSynthesis.speak(utterance);
    } else {
        console.log('浏览器不支持语音合成');
    }
}

// 开始学习（随机显示部分内容）
function startLesson(subject) {
    const content = document.getElementById(subject + '-content');
    const lessonsPerPage = 15; // 每次显示15个内容
    
    // 随机选择学习内容
    currentLessons = getRandomItems(learningData[subject].lessons, lessonsPerPage);
    
    let html = `
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px;">
            <h2 style="color: #667eea; font-size: 2em; margin: 0;">📚 学习内容</h2>
            <button class="subject-btn" onclick="startLesson('${subject}')" style="padding: 10px 25px;">
                🔄 换一批
            </button>
        </div>
    `;
    
    if (subject === 'english') {
        // 按分类组织英语单词
        const categories = {};
        currentLessons.forEach(lesson => {
            const cat = lesson.category || '其他';
            if (!categories[cat]) categories[cat] = [];
            categories[cat].push(lesson);
        });
        
        // 显示每个分类
        Object.keys(categories).forEach(category => {
            html += `<h3 style="color: #764ba2; margin: 30px 0 20px 0; font-size: 1.8em;">📌 ${category}</h3>`;
            categories[category].forEach(lesson => {
                html += `
                    <div class="lesson-item">
                        <div style="display: flex; justify-content: space-between; align-items: center;">
                            <div>
                                <h3>${lesson.emoji} ${lesson.word}</h3>
                                <p>中文：${lesson.translation}</p>
                            </div>
                            <button class="speak-btn" onclick="speak('${lesson.word}', 'en-US')" title="点击发音">
                                🔊
                            </button>
                        </div>
                    </div>
                `;
            });
        });
        
        html += `<p style="text-align: center; margin-top: 30px; color: #666; font-size: 1.2em;">
            本次显示 ${currentLessons.length} 个单词 | 题库共 ${learningData.english.lessons.length} 个单词
        </p>`;
        
    } else if (subject === 'math') {
        currentLessons.forEach(lesson => {
            html += `
                <div class="lesson-item">
                    <h3>${lesson.emoji || '📐'} ${lesson.title}</h3>
                    <p style="font-size: 1.4em;">${lesson.content}</p>
                </div>
            `;
        });
        
        html += `<p style="text-align: center; margin-top: 30px; color: #666; font-size: 1.2em;">
            本次显示 ${currentLessons.length} 个知识点 | 题库共 ${learningData.math.lessons.length} 个知识点
        </p>`;
        
    } else if (subject === 'chinese') {
        // 按分类组织汉字
        const categories = {};
        currentLessons.forEach(lesson => {
            const cat = lesson.category || '其他';
            if (!categories[cat]) categories[cat] = [];
            categories[cat].push(lesson);
        });
        
        // 显示每个分类
        Object.keys(categories).forEach(category => {
            html += `<h3 style="color: #764ba2; margin: 30px 0 20px 0; font-size: 1.8em;">📌 ${category}</h3>`;
            categories[category].forEach(lesson => {
                html += `
                    <div class="lesson-item">
                        <h3 style="font-size: 3em; color: #667eea;">${lesson.char}</h3>
                        <p style="font-size: 1.5em;">拼音：<span style="color: #e74c3c;">${lesson.pinyin}</span></p>
                        <p style="font-size: 1.3em;">意思：${lesson.meaning}</p>
                    </div>
                `;
            });
        });
        
        html += `<p style="text-align: center; margin-top: 30px; color: #666; font-size: 1.2em;">
            本次显示 ${currentLessons.length} 个汉字 | 题库共 ${learningData.chinese.lessons.length} 个汉字
        </p>`;
    }
    
    content.innerHTML = html;
}

// 开始测试（随机抽取题目）
function startTest(subject) {
    currentTest = subject;
    currentQuestionIndex = 0;
    score = 0;
    
    // 随机选择10道题
    const questionsPerTest = 10;
    currentTests = getRandomItems(learningData[subject].tests, questionsPerTest);
    
    showQuestion();
}

function showQuestion() {
    const content = document.getElementById(currentTest + '-content');
    
    if (currentQuestionIndex >= currentTests.length) {
        showScore();
        return;
    }
    
    const question = currentTests[currentQuestionIndex];
    let html = `
        <div class="test-question">
            <h3>问题 ${currentQuestionIndex + 1}/${currentTests.length}</h3>
    `;
    
    // 如果是英语测试，添加发音按钮
    if (currentTest === 'english') {
        // 提取英文单词（如果问题中包含英文单词）
        const wordMatch = question.question.match(/[A-Za-z]+/);
        if (wordMatch) {
            const word = wordMatch[0];
            html += `
                <div style="display: flex; align-items: center; justify-content: center; gap: 15px; margin: 20px 0;">
                    <p style="font-size: 1.8em; color: #333;">${question.question}</p>
                    <button class="speak-btn-small" onclick="speak('${word}', 'en-US')" title="点击发音">
                        🔊
                    </button>
                </div>
            `;
        } else {
            html += `<p style="font-size: 1.8em; margin: 20px 0; color: #333;">${question.question}</p>`;
        }
    } else {
        html += `<p style="font-size: 1.8em; margin: 20px 0; color: #333;">${question.question}</p>`;
    }
    
    html += `<div class="options">`;
    
    question.options.forEach((option, index) => {
        // 如果是英语测试且选项是英文，添加小的发音按钮
        if (currentTest === 'english' && /^[A-Za-z]+$/.test(option)) {
            html += `
                <div style="display: flex; align-items: center; gap: 10px;">
                    <button class="option-btn" onclick="checkAnswer(${index})" style="flex: 1;">${option}</button>
                    <button class="speak-btn-mini" onclick="event.stopPropagation(); speak('${option}', 'en-US')" title="发音">
                        🔊
                    </button>
                </div>
            `;
        } else {
            html += `
                <button class="option-btn" onclick="checkAnswer(${index})">${option}</button>
            `;
        }
    });
    
    html += `
            </div>
        </div>
    `;
    
    content.innerHTML = html;
}

function checkAnswer(selectedIndex) {
    const question = currentTests[currentQuestionIndex];
    const buttons = document.querySelectorAll('.option-btn');
    
    buttons.forEach((btn, index) => {
        btn.disabled = true;
        if (index === question.answer) {
            btn.classList.add('correct');
        } else if (index === selectedIndex && selectedIndex !== question.answer) {
            btn.classList.add('wrong');
        }
    });
    
    if (selectedIndex === question.answer) {
        score++;
        playSound('correct');
    } else {
        playSound('wrong');
    }
    
    setTimeout(() => {
        currentQuestionIndex++;
        showQuestion();
    }, 1500);
}

function showScore() {
    const content = document.getElementById(currentTest + '-content');
    const total = currentTests.length;
    const percentage = Math.round((score / total) * 100);
    
    let emoji = '🎉';
    let message = '太棒了！';
    
    if (percentage < 60) {
        emoji = '💪';
        message = '继续加油！';
    } else if (percentage < 80) {
        emoji = '👍';
        message = '做得不错！';
    }
    
    content.innerHTML = `
        <div class="score-display">
            <h2>${emoji} 测试完成！</h2>
            <p class="score">${score} / ${total}</p>
            <p style="font-size: 1.5em; margin: 20px 0;">正确率：${percentage}%</p>
            <p style="font-size: 1.3em; color: #666;">${message}</p>
            <button class="subject-btn" onclick="startTest('${currentTest}')" style="margin-top: 30px;">
                🔄 再测一次
            </button>
        </div>
    `;
}

function playSound(type) {
    // 简单的音效提示（可以替换为真实音频）
    if (type === 'correct') {
        console.log('✓ 正确！');
    } else {
        console.log('✗ 错误！');
    }
}

// 初始化
document.addEventListener('DOMContentLoaded', function() {
    console.log('幼儿园学习APP已启动！');
});
