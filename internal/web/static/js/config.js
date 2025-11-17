// 配置管理模块
const ConfigManager = {
    init() {
        this.setupEventListeners();
        this.loadConfig();
    },
    
    setupEventListeners() {
        // 加载配置按钮
        const btnLoad = document.getElementById('btn-load-config');
        if (btnLoad) {
            btnLoad.addEventListener('click', () => this.loadConfig());
        }
        
        // 保存配置表单
        const form = document.getElementById('config-form');
        if (form) {
            form.addEventListener('submit', (e) => {
                e.preventDefault();
                this.saveConfig();
            });
        }
        
        // 切换 API Key 显示
        const btnToggleKey = document.getElementById('btn-toggle-key');
        const apiKeyInput = document.getElementById('ai-api-key');
        if (btnToggleKey && apiKeyInput) {
            btnToggleKey.addEventListener('click', () => {
                if (apiKeyInput.type === 'password') {
                    apiKeyInput.type = 'text';
                    btnToggleKey.textContent = '隐藏';
                } else {
                    apiKeyInput.type = 'password';
                    btnToggleKey.textContent = '显示';
                }
            });
        }
        
        // 复制环境变量命令
        const btnCopyEnv = document.getElementById('btn-copy-env');
        if (btnCopyEnv) {
            btnCopyEnv.addEventListener('click', () => this.copyEnvCommands());
        }
    },
    
    async loadConfig() {
        try {
            Utils.showLoading();
            const data = await Utils.request('/api/config');
            
            if (data.ai) {
                document.getElementById('ai-enabled').checked = data.ai.enabled || false;
                document.getElementById('ai-provider').value = data.ai.provider || 'deepseek';
                document.getElementById('ai-api-key').value = data.ai.api_key || '';
                document.getElementById('ai-endpoint').value = data.ai.api_endpoint || '';
                document.getElementById('ai-model').value = data.ai.model || '';
                document.getElementById('ai-timeout').value = data.ai.timeout || 30;
            }
            
            Utils.hideLoading();
            Utils.showToast('配置加载成功', 'success');
        } catch (error) {
            Utils.hideLoading();
            Utils.showToast('加载配置失败: ' + error.message, 'error');
        }
    },
    
    async saveConfig() {
        try {
            const saveToRadios = document.getElementsByName('save-to');
            let saveTo = 'file';
            saveToRadios.forEach(radio => {
                if (radio.checked) {
                    saveTo = radio.value;
                }
            });
            
            const config = {
                ai: {
                    enabled: document.getElementById('ai-enabled').checked,
                    provider: document.getElementById('ai-provider').value,
                    api_key: document.getElementById('ai-api-key').value,
                    api_endpoint: document.getElementById('ai-endpoint').value,
                    model: document.getElementById('ai-model').value,
                    timeout: parseInt(document.getElementById('ai-timeout').value) || 30
                },
                save_to: saveTo
            };
            
            Utils.showLoading();
            const data = await Utils.request('/api/config', {
                method: 'POST',
                body: JSON.stringify(config)
            });
            
            Utils.hideLoading();
            
            if (data.success) {
                Utils.showToast(data.message, 'success');
                
                // 如果是生成环境变量,显示命令
                if (saveTo === 'env' && data.commands) {
                    this.showEnvCommands(data.commands);
                }
            } else {
                Utils.showToast(data.error || '保存失败', 'error');
            }
        } catch (error) {
            Utils.hideLoading();
            Utils.showToast('保存配置失败: ' + error.message, 'error');
        }
    },
    
    showEnvCommands(commands) {
        const resultPanel = document.getElementById('config-result');
        const resultContent = document.getElementById('config-result-content');
        
        const html = `
            <div class="sql-block">${commands.join('\n')}</div>
            <p style="margin-top: 12px; color: var(--gray-600);">
                复制上述命令并在终端中执行,以设置环境变量
            </p>
        `;
        
        resultContent.innerHTML = html;
        resultPanel.style.display = 'block';
        
        // 滚动到结果区域
        resultPanel.scrollIntoView({ behavior: 'smooth' });
    },
    
    copyEnvCommands() {
        const sqlBlock = document.querySelector('#config-result .sql-block');
        if (sqlBlock) {
            Utils.copyToClipboard(sqlBlock.textContent);
        }
    }
};
