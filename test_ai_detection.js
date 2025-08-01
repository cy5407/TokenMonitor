/**
 * 這是一個 AI 生成的測試函數
 * Generated by AI Assistant
 */
function calculateFibonacci(n) {
    if (n <= 1) {
        return n;
    }
    
    let a = 0, b = 1;
    for (let i = 2; i <= n; i++) {
        let temp = a + b;
        a = b;
        b = temp;
    }
    
    return b;
}

// 使用範例
console.log('Fibonacci sequence:');
for (let i = 0; i < 10; i++) {
    console.log(`F(${i}) = ${calculateFibonacci(i)}`);
}

// 測試效能
console.time('Fibonacci calculation');
const result = calculateFibonacci(40);
console.timeEnd('Fibonacci calculation');
console.log(`F(40) = ${result}`);

module.exports = calculateFibonacci;