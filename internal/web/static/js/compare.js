// 表结构比对模块
const CompareManager = {
    currentResult: null,
    
    init() {
        this.setupEventListeners();
        this.setupExamples();
    },
    
    setupEventListeners() {
        // 比对按钮
        const btnCompare = document.getElementById('btn-compare');
        if (btnCompare) {
            btnCompare.addEventListener('click', () => this.compare());
        }
        
        // 复制结果
        const btnCopy = document.getElementById('btn-copy-result');
        if (btnCopy) {
            btnCopy.addEventListener('click', () => this.copyResult());
        }
        
        // 下载结果
        const btnDownload = document.getElementById('btn-download-result');
        if (btnDownload) {
            btnDownload.addEventListener('click', () => this.downloadResult());
        }
    },
    
    // 设置示例按钮
    setupExamples() {
        const exampleSource = document.getElementById('example-source');
        const exampleTarget = document.getElementById('example-target');
        
        if (exampleSource) {
            exampleSource.addEventListener('click', () => {
                document.getElementById('source-sql').value = this.getExampleSourceSQL();
                document.getElementById('source-sql').dispatchEvent(new Event('input'));
            });
        }
        
        if (exampleTarget) {
            exampleTarget.addEventListener('click', () => {
                document.getElementById('target-sql').value = this.getExampleTargetSQL();
                document.getElementById('target-sql').dispatchEvent(new Event('input'));
            });
        }
    },
    
    // 获取示例源表SQL
    getExampleSourceSQL() {
        return `CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100),
    created_at DATETIME DEFAULT NOW()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`;
    },
    
    // 获取示例目标表SQL
    getExampleTargetSQL() {
        return `CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    status TINYINT DEFAULT 1,
    created_at DATETIME DEFAULT NOW(),
    updated_at DATETIME
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`;
    },
    
    async compare() {
        const sourceSQL = document.getElementById('source-sql').value.trim();
        const targetSQL = document.getElementById('target-sql').value.trim();
        const enableAI = document.getElementById('enable-ai-compare').checked;
        
        if (!sourceSQL || !targetSQL) {
            Utils.showToast('请输入源表和目标表的 SQL 语句', 'error');
            return;
        }
        
        try {
            Utils.showLoading();
            const data = await Utils.request('/api/compare', {
                method: 'POST',
                body: JSON.stringify({
                    source_sql: sourceSQL,
                    target_sql: targetSQL,
                    enable_ai: enableAI
                })
            });
            
            Utils.hideLoading();
            
            if (data.success === false) {
                Utils.showToast(data.error || '比对失败', 'error');
                return;
            }
            
            this.currentResult = data;
            this.displayResult(data);
            Utils.showToast('比对完成', 'success');
        } catch (error) {
            Utils.hideLoading();
            Utils.showToast('比对失败: ' + error.message, 'error');
        }
    },
    
    displayResult(data) {
        const resultPanel = document.getElementById('compare-result');
        const resultContent = document.getElementById('compare-result-content');
        
        if (!data.has_changes) {
            resultContent.innerHTML = `
                <div class="alert alert-success">
                    <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                    <span>两个表结构完全相同，无需修改！</span>
                </div>
            `;
            resultPanel.style.display = 'block';
            return;
        }
        
        let html = `
            <div class="alert alert-info mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                <div>
                    <div class="font-bold">差异摘要</div>
                    <div class="text-sm whitespace-pre-wrap">${data.summary}</div>
                </div>
            </div>
        `;
        
        if (data.ddls && data.ddls.length > 0) {
            html += '<div class="mb-4"><h4 class="font-semibold mb-3 text-lg">DDL 语句</h4>';
            html += '<div class="space-y-2">';
            data.ddls.forEach(ddl => {
                let badgeClass = 'badge-success';
                let iconHtml = '<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>';
                
                if (ddl.includes('ADD COLUMN')) {
                    badgeClass = 'badge-success';
                } else if (ddl.includes('MODIFY COLUMN')) {
                    badgeClass = 'badge-warning';
                    iconHtml = '<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>';
                } else if (ddl.includes('DROP')) {
                    badgeClass = 'badge-error';
                    iconHtml = '<svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="5" y1="12" x2="19" y2="12"></line></svg>';
                }
                
                html += `
                    <div class="bg-base-200 rounded-lg p-4 hover:bg-base-300 transition-colors">
                        <div class="flex items-start gap-3">
                            <span class="badge ${badgeClass} gap-1">
                                ${iconHtml}
                            </span>
                            <code class="flex-1 text-sm font-mono">${ddl};</code>
                        </div>
                    </div>
                `;
            });
            html += '</div></div>';
        }
        
        if (data.ai_analysis) {
            html += '<div class="divider"></div>';
            html += '<div class="bg-gradient-to-r from-primary/10 to-secondary/10 rounded-lg p-6">';
            html += '<h4 class="font-bold text-lg mb-3 flex items-center gap-2">';
            html += '<svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 3c-1.2 0-2.4.6-3 1.7A3.6 3.6 0 0 0 4.6 9c-1 .6-1.7 1.8-1.7 3s.7 2.4 1.7 3c-.3 1.2 0 2.5 1 3.4.8.8 2.1 1.2 3.3 1 .6 1 1.8 1.6 3 1.6s2.4-.6 3-1.7c1.2.3 2.5 0 3.4-1 .8-.8 1.2-2.1 1-3.3 1-.6 1.6-1.8 1.6-3s-.6-2.4-1.7-3c.3-1.2 0-2.5-1-3.4a3.7 3.7 0 0 0-3.3-1c-.6-1-1.8-1.6-3-1.6Z"></path><path d="m9 12 2 2 4-4"></path></svg>';
            html += 'AI 分析结果</h4>';
            
            if (data.ai_analysis.summary) {
                html += `<p class="text-base-content/80 mb-4">${data.ai_analysis.summary}</p>`;
            }
            
            if (data.ai_analysis.suggestions && data.ai_analysis.suggestions.length > 0) {
                html += '<div class="collapse collapse-arrow bg-base-100 mb-2">';
                html += '<input type="checkbox" checked />';
                html += '<div class="collapse-title font-medium">✨ 优化建议</div>';
                html += '<div class="collapse-content"><ul class="list-disc list-inside space-y-1">';
                data.ai_analysis.suggestions.forEach(s => {
                    html += `<li class="text-sm">${s}</li>`;
                });
                html += '</ul></div></div>';
            }
            
            if (data.ai_analysis.risks && data.ai_analysis.risks.length > 0) {
                html += '<div class="collapse collapse-arrow bg-base-100 mb-2">';
                html += '<input type="checkbox" />';
                html += '<div class="collapse-title font-medium">⚠️ 潜在风险</div>';
                html += '<div class="collapse-content"><ul class="list-disc list-inside space-y-1">';
                data.ai_analysis.risks.forEach(r => {
                    html += `<li class="text-sm text-warning">${r}</li>`;
                });
                html += '</ul></div></div>';
            }
            
            if (data.ai_analysis.best_practice && data.ai_analysis.best_practice.length > 0) {
                html += '<div class="collapse collapse-arrow bg-base-100">';
                html += '<input type="checkbox" />';
                html += '<div class="collapse-title font-medium">📚 最佳实践</div>';
                html += '<div class="collapse-content"><ul class="list-disc list-inside space-y-1">';
                data.ai_analysis.best_practice.forEach(bp => {
                    html += `<li class="text-sm">${bp}</li>`;
                });
                html += '</ul></div></div>';
            }
            
            html += '</div>';
        }
        
        resultContent.innerHTML = html;
        resultPanel.style.display = 'block';
        
        // 滚动到结果区域
        resultPanel.scrollIntoView({ behavior: 'smooth' });
    },
    
    copyResult() {
        if (!this.currentResult || !this.currentResult.ddls) {
            Utils.showToast('没有可复制的内容', 'error');
            return;
        }
        
        const text = this.currentResult.ddls.map(ddl => ddl + ';').join('\n');
        Utils.copyToClipboard(text);
    },
    
    downloadResult() {
        if (!this.currentResult || !this.currentResult.ddls) {
            Utils.showToast('没有可下载的内容', 'error');
            return;
        }
        
        const content = this.currentResult.ddls.map(ddl => ddl + ';').join('\n');
        Utils.downloadFile(content, 'migration.sql');
    }
};
