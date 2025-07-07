// 用户搜索功能
document.addEventListener('DOMContentLoaded', function() {
    // 添加页面加载动画
    document.body.style.opacity = '0';
    setTimeout(() => {
        document.body.style.transition = 'opacity 0.5s ease-out';
        document.body.style.opacity = '1';
    }, 100);

    // 搜索输入框
    const searchInput = document.getElementById('searchInput');
    const filterRole = document.getElementById('filterRole');
    const usersTableBody = document.getElementById('usersTableBody');
    const emptyState = document.getElementById('emptyState');

    // 搜索和筛选功能
    function filterUsers() {
        if (!searchInput || !usersTableBody) return;

        const searchTerm = searchInput.value.toLowerCase();
        const roleFilter = filterRole ? filterRole.value : 'all';
        const rows = usersTableBody.getElementsByClassName('user-row');
        let visibleCount = 0;

        for (let row of rows) {
            const username = row.querySelector('.user-info span:last-child').textContent.toLowerCase();
            const email = row.cells[2].textContent.toLowerCase();
            const role = row.getAttribute('data-role');

            const matchesSearch = username.includes(searchTerm) || email.includes(searchTerm);
            const matchesRole = roleFilter === 'all' || role === roleFilter;

            if (matchesSearch && matchesRole) {
                row.style.display = '';
                row.style.animation = 'fadeIn 0.3s ease-out';
                visibleCount++;
            } else {
                row.style.display = 'none';
            }
        }

        // 显示空状态
        if (emptyState) {
            emptyState.style.display = visibleCount === 0 ? 'block' : 'none';
        }
    }

    // 添加事件监听器
    if (searchInput) {
        searchInput.addEventListener('input', filterUsers);

        // 添加搜索框聚焦动画
        searchInput.addEventListener('focus', function() {
            this.parentElement.style.transform = 'scale(1.02)';
        });

        searchInput.addEventListener('blur', function() {
            this.parentElement.style.transform = 'scale(1)';
        });
    }

    if (filterRole) {
        filterRole.addEventListener('change', filterUsers);
    }

    // 密码强度检测
    const passwordInput = document.getElementById('password');
    if (passwordInput && document.querySelector('.password-strength')) {
        passwordInput.addEventListener('input', function() {
            const password = this.value;
            const strengthFill = document.querySelector('.strength-fill');

            if (!strengthFill) return;

            let strength = 0;
            if (password.length >= 6) strength++;
            if (password.length >= 10) strength++;
            if (/[A-Z]/.test(password) && /[a-z]/.test(password)) strength++;
            if (/[0-9]/.test(password)) strength++;
            if (/[^A-Za-z0-9]/.test(password)) strength++;

            strengthFill.className = 'strength-fill';
            if (strength <= 2) {
                strengthFill.classList.add('weak');
            } else if (strength <= 3) {
                strengthFill.classList.add('medium');
            } else {
                strengthFill.classList.add('strong');
            }
        });
    }

    // 表单提交处理
    const forms = document.querySelectorAll('form');
    forms.forEach(form => {
        // 登录表单特殊处理
        if (form.action && form.action.includes('/login')) {
            form.addEventListener('submit', function(e) {
                const submitBtn = form.querySelector('button[type="submit"]');
                if (submitBtn) {
                    submitBtn.disabled = true;
                    submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> 处理中...';
                    submitBtn.style.transform = 'scale(0.95)';
                }
            });
        }

        // 注册表单特殊处理
        if (form.action && form.action.includes('/register')) {
            form.addEventListener('submit', function(e) {
                const submitBtn = form.querySelector('button[type="submit"]');
                if (submitBtn) {
                    submitBtn.disabled = true;
                    submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> 创建账户...';
                    submitBtn.style.transform = 'scale(0.95)';
                }
            });
        }
    });

    // 添加输入框动画效果
    const inputs = document.querySelectorAll('input, select');
    inputs.forEach(input => {
        input.addEventListener('focus', function() {
            this.style.transform = 'scale(1.02)';
        });

        input.addEventListener('blur', function() {
            this.style.transform = 'scale(1)';
        });
    });

    // 点击模态框外部关闭
    window.addEventListener('click', function(e) {
        const modal = document.getElementById('editModal');
        if (modal && e.target === modal) {
            closeModal();
        }
    });

    // 添加按钮悬停效果
    const buttons = document.querySelectorAll('.btn-primary, .btn-secondary, .btn-icon');
    buttons.forEach(button => {
        button.addEventListener('mouseenter', function() {
            this.style.transform = 'translateY(-2px)';
        });

        button.addEventListener('mouseleave', function() {
            this.style.transform = 'translateY(0)';
        });
    });

    // 统计管理员数量（如果在用户页面）
    const adminCountEl = document.getElementById('adminCount');
    if (adminCountEl) {
        const adminCount = document.querySelectorAll('.user-row[data-role="admin"]').length;
        adminCountEl.textContent = adminCount;

        // 添加数字动画
        animateValue(adminCountEl, 0, adminCount, 1000);
    }

    // 统计总用户数动画
    const statValues = document.querySelectorAll('.stat-value');
    statValues.forEach(stat => {
        const finalValue = parseInt(stat.textContent);
        if (!isNaN(finalValue)) {
            animateValue(stat, 0, finalValue, 1000);
        }
    });
});

// 数字动画函数
function animateValue(obj, start, end, duration) {
    let startTimestamp = null;
    const step = (timestamp) => {
        if (!startTimestamp) startTimestamp = timestamp;
        const progress = Math.min((timestamp - startTimestamp) / duration, 1);
        obj.innerHTML = Math.floor(progress * (end - start) + start);
        if (progress < 1) {
            window.requestAnimationFrame(step);
        }
    };
    window.requestAnimationFrame(step);
}

// 全局函数（供HTML内联调用）
function editUser(id, username, email, role) {
    const modal = document.getElementById('editModal');

    document.getElementById('edit-user-id').value = id;
    document.getElementById('edit-username').value = username;
    document.getElementById('edit-email').value = email;
    document.getElementById('edit-role').value = role;

    // 显示模态框并添加动画
    modal.style.display = 'flex';
    modal.style.opacity = '0';
    setTimeout(() => {
        modal.style.transition = 'opacity 0.3s ease-out';
        modal.style.opacity = '1';
    }, 10);
}

function closeModal() {
    const modal = document.getElementById('editModal');
    modal.style.opacity = '0';
    setTimeout(() => {
        modal.style.display = 'none';
        modal.style.opacity = '1';
    }, 300);
}

function confirmDelete(username) {
    // 创建自定义确认对话框
    const confirmed = confirm(`确定要删除用户 "${username}" 吗？此操作不可撤销。`);

    if (confirmed) {
        // 添加删除动画
        event.target.closest('.user-row').style.animation = 'fadeOut 0.3s ease-out';
    }

    return confirmed;
}

// 添加页面滚动效果
window.addEventListener('scroll', function() {
    const navbar = document.querySelector('.navbar');
    if (navbar) {
        if (window.scrollY > 50) {
            navbar.style.background = 'rgba(255, 255, 255, 0.1)';
        } else {
            navbar.style.background = 'rgba(255, 255, 255, 0.05)';
        }
    }
});