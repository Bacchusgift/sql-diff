// 表结构预览模块
const TablePreview = {
    debounceTimers: {},
    
    // 初始化预览功能
    init() {
        this.setupPreview('source-sql', 'source-preview');
        this.setupPreview('target-sql', 'target-preview');
    },
    
    // 设置预览监听
    setupPreview(textareaId, previewId) {
        const textarea = document.getElementById(textareaId);
        if (!textarea) return;
        
        textarea.addEventListener('input', () => {
            // 防抖处理
            if (this.debounceTimers[textareaId]) {
                clearTimeout(this.debounceTimers[textareaId]);
            }
            
            this.debounceTimers[textareaId] = setTimeout(() => {
                this.parseAndDisplay(textarea.value, previewId);
            }, 500);
        });
    },
    
    // 解析并显示
    async parseAndDisplay(sql, previewId) {
        const previewDiv = document.getElementById(previewId);
        if (!previewDiv) return;
        
        const trimmedSQL = sql.trim();
        if (!trimmedSQL) {
            previewDiv.style.display = 'none';
            return;
        }
        
        try {
            const data = await Utils.request('/api/parse-create', {
                method: 'POST',
                body: JSON.stringify({ sql: trimmedSQL })
            });
            
            if (data.success) {
                this.renderTable(data, previewDiv);
                previewDiv.style.display = 'block';
            } else {
                this.renderError(data.error, previewDiv);
                previewDiv.style.display = 'block';
            }
        } catch (error) {
            console.error('解析失败:', error);
            previewDiv.style.display = 'none';
        }
    },
    
    // 渲染表格
    renderTable(data, container) {
        let html = `
            <div class="divider">表结构预览</div>
            <div class="bg-base-200 rounded-lg p-4">
                <div class="flex items-center gap-2 mb-3">
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M3 6h18"></path>
                        <path d="M3 12h18"></path>
                        <path d="M3 18h18"></path>
                    </svg>
                    <span class="font-bold text-lg">${data.table_name}</span>
                </div>
                <div class="overflow-x-auto">
                    <table class="table table-sm table-zebra">
                        <thead>
                            <tr>
                                <th>列名</th>
                                <th>类型</th>
                                <th>长度</th>
                                <th>默认值</th>
                                <th>NULL</th>
                                <th>约束</th>
                                <th>备注</th>
                            </tr>
                        </thead>
                        <tbody>
        `;
        
        data.columns.forEach(col => {
            const constraints = [];
            if (col.is_primary) constraints.push('<span class="badge badge-sm badge-primary">PK</span>');
            if (col.is_unique) constraints.push('<span class="badge badge-sm badge-info">UQ</span>');
            if (col.auto_inc) constraints.push('<span class="badge badge-sm badge-success">AI</span>');
            if (col.unsigned) constraints.push('<span class="badge badge-sm badge-ghost">UN</span>');
            
            html += `
                <tr>
                    <td class="font-mono font-semibold">${col.name}</td>
                    <td class="font-mono text-primary">${col.type}</td>
                    <td class="text-center">${col.length || '-'}</td>
                    <td class="font-mono text-sm">${col.default || '-'}</td>
                    <td class="text-center">${col.nullable ? '<span class="text-success">✓</span>' : '<span class="text-error">✗</span>'}</td>
                    <td>${constraints.join(' ') || '-'}</td>
                    <td class="text-sm text-base-content/60">${col.comment || '-'}</td>
                </tr>
            `;
        });
        
        html += `
                        </tbody>
                    </table>
                </div>
                <div class="text-xs text-base-content/50 mt-2">
                    共 ${data.columns.length} 个字段
                </div>
            </div>
        `;
        
        container.innerHTML = html;
    },
    
    // 渲染错误
    renderError(error, container) {
        container.innerHTML = `
            <div class="alert alert-error mt-4">
                <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span>${error}</span>
            </div>
        `;
    }
};

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', () => {
    TablePreview.init();
});
