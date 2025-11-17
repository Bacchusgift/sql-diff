// 主题管理模块
const ThemeManager = {
    currentTheme: 'light',
    
    init() {
        this.loadTheme();
        this.setupToggle();
    },
    
    // 加载主题
    loadTheme() {
        const savedTheme = localStorage.getItem('sql-diff-theme');
        if (savedTheme) {
            this.currentTheme = savedTheme;
        } else {
            // 检测系统偏好
            if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
                this.currentTheme = 'dark';
            }
        }
        this.applyTheme();
    },
    
    // 应用主题
    applyTheme() {
        const html = document.documentElement;
        html.setAttribute('data-theme', this.currentTheme);
        
        // 更新图标
        const icon = document.querySelector('#theme-toggle i');
        if (icon) {
            icon.setAttribute('data-lucide', this.currentTheme === 'dark' ? 'moon' : 'sun');
            lucide.createIcons();
        }
    },
    
    // 切换主题
    toggleTheme() {
        this.currentTheme = this.currentTheme === 'light' ? 'dark' : 'light';
        this.applyTheme();
        localStorage.setItem('sql-diff-theme', this.currentTheme);
    },
    
    // 设置切换按钮
    setupToggle() {
        const toggleBtn = document.getElementById('theme-toggle');
        if (toggleBtn) {
            toggleBtn.addEventListener('click', () => this.toggleTheme());
        }
    }
};

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', () => {
    ThemeManager.init();
});
