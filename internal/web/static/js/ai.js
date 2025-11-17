// AI 功能模块
const AIManager = {
    currentCreateResult: null,
    currentAlterResult: null,
    
    init() {
        this.setupEventListeners();
    },
    
    setupEventListeners() {
        // AI 生成 CREATE TABLE
        const btnGenerateCreate = document.getElementById('btn-generate-create');
        if (btnGenerateCreate) {
            btnGenerateCreate.addEventListener('click', () => this.generateCreate());
        }
        
        const btnCopyCreate = document.getElementById('btn-copy-create');
        if (btnCopyCreate) {
            btnCopyCreate.addEventListener('click', () => this.copyCreate());
        }
        
        const btnDownloadCreate = document.getElementById('btn-download-create');
        if (btnDownloadCreate) {
            btnDownloadCreate.addEventListener('click', () => this.downloadCreate());
        }
        
        // AI 生成 ALTER TABLE
        const btnGenerateAlter = document.getElementById('btn-generate-alter');
        if (btnGenerateAlter) {
            btnGenerateAlter.addEventListener('click', () => this.generateAlter());
        }
        
        const btnCopyAlter = document.getElementById('btn-copy-alter');
        if (btnCopyAlter) {
            btnCopyAlter.addEventListener('click', () => this.copyAlter());
        }
        
        const btnDownloadAlter = document.getElementById('btn-download-alter');
        if (btnDownloadAlter) {
            btnDownloadAlter.addEventListener('click', () => this.downloadAlter());
        }
    },
    
    async generateCreate() {
        const description = document.getElementById('create-description').value.trim();
        
        if (!description) {
            Utils.showToast('请输入表结构描述', 'error');
            return;
        }
        
        try {
            Utils.showLoading();
            const data = await Utils.request('/api/ai/generate-create', {
                method: 'POST',
                body: JSON.stringify({ description })
            });
            
            Utils.hideLoading();
            
            if (!data.success) {
                Utils.showToast(data.error || 'AI 生成失败', 'error');
                return;
            }
            
            this.currentCreateResult = data.sql;
            this.displayCreateResult(data.sql);
            Utils.showToast('生成成功', 'success');
        } catch (error) {
            Utils.hideLoading();
            Utils.showToast('生成失败: ' + error.message, 'error');
        }
    },
    
    displayCreateResult(sql) {
        const resultPanel = document.getElementById('create-result');
        const resultContent = document.getElementById('create-result-content');
        
        resultContent.innerHTML = `
            <div class="sql-block">${sql};</div>
        `;
        
        resultPanel.style.display = 'block';
        resultPanel.scrollIntoView({ behavior: 'smooth' });
    },
    
    copyCreate() {
        if (!this.currentCreateResult) {
            Utils.showToast('没有可复制的内容', 'error');
            return;
        }
        
        Utils.copyToClipboard(this.currentCreateResult + ';');
    },
    
    downloadCreate() {
        if (!this.currentCreateResult) {
            Utils.showToast('没有可下载的内容', 'error');
            return;
        }
        
        Utils.downloadFile(this.currentCreateResult + ';', 'create_table.sql');
    },
    
    async generateAlter() {
        const currentDDL = document.getElementById('current-ddl').value.trim();
        const description = document.getElementById('alter-description').value.trim();
        
        if (!currentDDL || !description) {
            Utils.showToast('请输入表结构和修改描述', 'error');
            return;
        }
        
        try {
            Utils.showLoading();
            const data = await Utils.request('/api/ai/generate-alter', {
                method: 'POST',
                body: JSON.stringify({
                    current_ddl: currentDDL,
                    description: description
                })
            });
            
            Utils.hideLoading();
            
            if (!data.success) {
                Utils.showToast(data.error || 'AI 生成失败', 'error');
                return;
            }
            
            this.currentAlterResult = data.sqls;
            this.displayAlterResult(data.sqls);
            Utils.showToast('生成成功', 'success');
        } catch (error) {
            Utils.hideLoading();
            Utils.showToast('生成失败: ' + error.message, 'error');
        }
    },
    
    displayAlterResult(sqls) {
        const resultPanel = document.getElementById('alter-result');
        const resultContent = document.getElementById('alter-result-content');
        
        let html = '<ul class="ddl-list">';
        sqls.forEach(sql => {
            html += `<li class="ddl-item add">${sql};</li>`;
        });
        html += '</ul>';
        
        resultContent.innerHTML = html;
        resultPanel.style.display = 'block';
        resultPanel.scrollIntoView({ behavior: 'smooth' });
    },
    
    copyAlter() {
        if (!this.currentAlterResult) {
            Utils.showToast('没有可复制的内容', 'error');
            return;
        }
        
        const text = this.currentAlterResult.map(sql => sql + ';').join('\n');
        Utils.copyToClipboard(text);
    },
    
    downloadAlter() {
        if (!this.currentAlterResult) {
            Utils.showToast('没有可下载的内容', 'error');
            return;
        }
        
        const content = this.currentAlterResult.map(sql => sql + ';').join('\n');
        Utils.downloadFile(content, 'alter_table.sql');
    }
};
