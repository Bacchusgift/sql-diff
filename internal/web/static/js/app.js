// 全局应用对象
const App = {
    currentPage: 'compare',
    
    // 初始化
    init() {
        this.setupNavigation();
        this.setupRouting();
        this.loadCurrentPage();
    },
    
    // 设置导航
    setupNavigation() {
        const navButtons = document.querySelectorAll('.tab');
        navButtons.forEach(btn => {
            btn.addEventListener('click', () => {
                const page = btn.dataset.page;
                this.navigateTo(page);
            });
        });
    },
    
    // 设置路由
    setupRouting() {
        // 监听 hash 变化
        window.addEventListener('hashchange', () => {
            this.loadCurrentPage();
        });
    },
    
    // 导航到页面
    navigateTo(page) {
        window.location.hash = page;
    },
    
    // 加载当前页面
    loadCurrentPage() {
        const hash = window.location.hash.slice(1);
        const page = hash || 'compare';
        
        // 隐藏所有页面
        document.querySelectorAll('.page').forEach(p => {
            p.style.display = 'none';
        });
        
        // 显示当前页面
        const currentPage = document.getElementById('page-' + page);
        if (currentPage) {
            currentPage.style.display = 'block';
            this.currentPage = page;
        }
        
        // 更新导航按钮状态
        document.querySelectorAll('.tab').forEach(btn => {
            if (btn.dataset.page === page) {
                btn.classList.add('tab-active');
            } else {
                btn.classList.remove('tab-active');
            }
        });
        
        // 重新创建图标
        if (typeof lucide !== 'undefined') {
            lucide.createIcons();
        }
    }
};

// 工具函数
const Utils = {
    // 显示加载动画
    showLoading() {
        const loading = document.getElementById('loading');
        if (loading) {
            loading.classList.remove('hidden');
            loading.classList.add('flex');
        }
    },
    
    // 隐藏加载动画
    hideLoading() {
        const loading = document.getElementById('loading');
        if (loading) {
            loading.classList.add('hidden');
            loading.classList.remove('flex');
        }
    },
    
    // 显示 Toast 提示
    showToast(message, type = 'info') {
        const toast = document.getElementById('toast');
        const alertClass = {
            'success': 'alert-success',
            'error': 'alert-error',
            'warning': 'alert-warning',
            'info': 'alert-info'
        }[type] || 'alert-info';
        
        const icon = {
            'success': '<svg xmlns="http://www.w3.org/2000/svg" class="shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="white" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>',
            'error': '<svg xmlns="http://www.w3.org/2000/svg" class="shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="white" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>',
            'warning': '<svg xmlns="http://www.w3.org/2000/svg" class="shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="white" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>',
            'info': '<svg xmlns="http://www.w3.org/2000/svg" class="shrink-0 w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="white" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>'
        }[type] || '';
        
        toast.innerHTML = `<div class="alert ${alertClass} shadow-lg">${icon}<span>${message}</span></div>`;
        toast.classList.remove('hidden');
        
        setTimeout(() => {
            toast.classList.add('hidden');
        }, 3000);
    },
    
    // 复制到剪贴板
    async copyToClipboard(text) {
        try {
            await navigator.clipboard.writeText(text);
            this.showToast('已复制到剪贴板', 'success');
        } catch (err) {
            this.showToast('复制失败', 'error');
        }
    },
    
    // 下载文件
    downloadFile(content, filename) {
        const blob = new Blob([content], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = filename;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
        this.showToast('文件已下载', 'success');
    },
    
    // API 请求
    async request(url, options = {}) {
        try {
            const response = await fetch(url, {
                ...options,
                headers: {
                    'Content-Type': 'application/json',
                    ...options.headers
                }
            });
            
            const data = await response.json();
            return data;
        } catch (error) {
            throw new Error('网络请求失败: ' + error.message);
        }
    }
};

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', () => {
    App.init();
    
    // 初始化各个模块
    if (typeof ConfigManager !== 'undefined') {
        ConfigManager.init();
    }
    if (typeof CompareManager !== 'undefined') {
        CompareManager.init();
    }
    if (typeof AIManager !== 'undefined') {
        AIManager.init();
    }
});
