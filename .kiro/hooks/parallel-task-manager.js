// parallel-task-manager.js
// 併行任務管理器的執行邏輯

const tasks = [
  {
    id: "csv-report-generator",
    name: "CSV 格式報告生成器",
    priority: "high",
    estimatedTime: "2h",
    dependencies: ["basic-report-system"],
    files: [
      "internal/reporter/csv_generator.go",
      "internal/reporter/csv_generator_test.go"
    ]
  },
  {
    id: "html-report-generator", 
    name: "HTML 格式報告生成器",
    priority: "high",
    estimatedTime: "3h",
    dependencies: ["basic-report-system"],
    files: [
      "internal/reporter/html_generator.go",
      "internal/reporter/html_generator_test.go",
      "internal/reporter/templates/html_template.html"
    ]
  },
  {
    id: "data-persistence",
    name: "資料持久化系統",
    priority: "medium", 
    estimatedTime: "4h",
    dependencies: [],
    files: [
      "internal/storage/storage.go",
      "internal/storage/storage_test.go",
      "internal/storage/json_storage.go"
    ]
  },
  {
    id: "cli-enhancement",
    name: "CLI 介面增強",
    priority: "low",
    estimatedTime: "2h",
    dependencies: ["data-persistence"],
    files: [
      "cmd/enhanced_cli.go",
      "cmd/interactive.go"
    ]
  }
];

// 分析任務依賴關係
function analyzeDependencies() {
  const independentTasks = tasks.filter(task => task.dependencies.length === 0);
  const dependentTasks = tasks.filter(task => task.dependencies.length > 0);
  
  return {
    canStartImmediately: independentTasks,
    needsWaiting: dependentTasks
  };
}

// 建議併行執行策略
function suggestParallelStrategy() {
  const analysis = analyzeDependencies();
  
  return {
    phase1: analysis.canStartImmediately,
    phase2: analysis.needsWaiting,
    recommendation: "建議先併行執行 CSV 和 HTML 報告生成器，同時開始資料持久化系統的設計"
  };
}

// 生成任務執行報告
function generateTaskReport() {
  const strategy = suggestParallelStrategy();
  
  console.log("=== 併行任務執行計畫 ===");
  console.log("\n階段 1 (可立即開始):");
  strategy.phase1.forEach(task => {
    console.log(`- ${task.name} (${task.estimatedTime})`);
  });
  
  console.log("\n階段 2 (需等待依賴):");
  strategy.phase2.forEach(task => {
    console.log(`- ${task.name} (${task.estimatedTime}) - 依賴: ${task.dependencies.join(', ')}`);
  });
  
  console.log(`\n建議: ${strategy.recommendation}`);
}

// 執行入口
if (typeof module !== 'undefined' && module.exports) {
  module.exports = {
    tasks,
    analyzeDependencies,
    suggestParallelStrategy,
    generateTaskReport
  };
} else {
  // 在 Kiro 環境中執行
  generateTaskReport();
}