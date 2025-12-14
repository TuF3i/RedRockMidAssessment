// 全局变量
const API_BASE_URL = 'http://localhost:8080';
let currentUser = null;
let authToken = null;
let currentSection = 'selectable-courses';
let confirmCallback = null;

// 工具函数
function showLoading(button) {
    const btnText = button.querySelector('.btn-text');
    const btnLoading = button.querySelector('.btn-loading');
    if (btnText) btnText.style.display = 'none';
    if (btnLoading) btnLoading.style.display = 'inline';
    button.disabled = true;
}

function hideLoading(button) {
    const btnText = button.querySelector('.btn-text');
    const btnLoading = button.querySelector('.btn-loading');
    if (btnText) btnText.style.display = 'inline';
    if (btnLoading) btnLoading.style.display = 'none';
    button.disabled = false;
}

function showErrorModal(status, info) {
    document.getElementById('errorStatus').textContent = status || '未知';
    document.getElementById('errorInfo').textContent = info || '发生未知错误';
    document.getElementById('errorModal').style.display = 'block';
}

function closeErrorModal() {
    document.getElementById('errorModal').style.display = 'none';
}

function showConfirmModal(message, callback) {
    document.getElementById('confirmMessage').textContent = message;
    document.getElementById('confirmModal').style.display = 'block';
    confirmCallback = callback;
}

function closeConfirmModal() {
    document.getElementById('confirmModal').style.display = 'none';
    confirmCallback = null;
}

function confirmAction() {
    if (confirmCallback) {
        confirmCallback();
    }
    closeConfirmModal();
}

// API请求函数
async function apiRequest(url, options = {}) {
    const defaultOptions = {
        headers: {
            'Content-Type': 'application/json',
        }
    };

    if (authToken) {
        defaultOptions.headers['Authorization'] = `Bearer ${authToken}`;
    }

    const mergedOptions = {
        ...defaultOptions,
        ...options,
        headers: {
            ...defaultOptions.headers,
            ...options.headers
        }
    };

    try {
        const response = await fetch(`${API_BASE_URL}${url}`, mergedOptions);
        const data = await response.json();
        
        if (!response.ok || (data.status && data.status !== 20000)) {
            throw new Error(JSON.stringify({
                status: data.status || response.status,
                info: data.info || '请求失败'
            }));
        }
        
        return data;
    } catch (error) {
        if (error.message.includes('Failed to fetch')) {
            throw new Error(JSON.stringify({
                status: '网络错误',
                info: '无法连接到服务器，请检查网络连接'
            }));
        }
        throw error;
    }
}

// 数据校验函数
function validateStudentID(stuID) {
    if (!stuID || stuID.length !== 10) {
        return '学号必须是10位数字';
    }
    if (!/^\d{10}$/.test(stuID)) {
        return '学号只能包含数字';
    }
    return '';
}

function validateName(name) {
    if (!name || name.length < 2 || name.length > 15) {
        return '姓名长度必须在2-15个字符之间';
    }
    const reservedNames = ['admin', 'root', 'user', 'api', 'bob'];
    if (reservedNames.includes(name.toLowerCase())) {
        return '该用户名已被系统保留';
    }
    if (!/^[a-zA-Z0-9\u4e00-\u9fa5]([a-zA-Z0-9_\u4e00-\u9fa5]*[a-zA-Z0-9\u4e00-\u9fa5])?$/.test(name)) {
        return '姓名只能包含字母、数字、中文和下划线，且不能以下划线开头或结尾';
    }
    return '';
}

function validateClass(stuClass) {
    if (!stuClass || stuClass.length < 3 || stuClass.length > 15) {
        return '班级长度必须在3-15个字符之间';
    }
    return '';
}

function validatePassword(password) {
    if (!password || password.length < 6) {
        return '密码长度至少为6位';
    }
    return '';
}

function validateAge(age) {
    if (!age || age < 16 || age > 30) {
        return '年龄必须在16-30岁之间';
    }
    return '';
}

// 登录和注册功能
async function login() {
    const stuId = document.getElementById('stu_id').value.trim();
    const password = document.getElementById('password').value;
    const role = document.querySelector('input[name="role"]:checked').value;
    
    // 清除之前的错误信息
    document.getElementById('stu_id_error').textContent = '';
    document.getElementById('password_error').textContent = '';
    
    // 数据校验
    const stuIdError = validateStudentID(stuId);
    if (stuIdError) {
        document.getElementById('stu_id_error').textContent = stuIdError;
        return;
    }
    
    if (!password) {
        document.getElementById('password_error').textContent = '请输入密码';
        return;
    }
    
    const loginBtn = document.querySelector('.login-btn');
    showLoading(loginBtn);
    
    try {
        const response = await apiRequest('/v1/api/public/login', {
            method: 'POST',
            body: JSON.stringify({
                stu_id: stuId,
                password: password
            })
        });
        
        authToken = response.data.access_token;
        localStorage.setItem('authToken', authToken);
        localStorage.setItem('userRole', role);
        
        // 根据角色跳转到相应页面
        if (role === 'admin') {
            window.location.href = 'admin.html';
        } else {
            window.location.href = 'student.html';
        }
        
    } catch (error) {
        try {
            const errorData = JSON.parse(error.message);
            showErrorModal(errorData.status, errorData.info);
        } catch {
            showErrorModal('错误', error.message);
        }
    } finally {
        hideLoading(loginBtn);
    }
}

async function register() {
    const name = document.getElementById('reg_name').value.trim();
    const stuId = document.getElementById('reg_stu_id').value.trim();
    const stuClass = document.getElementById('reg_stu_class').value.trim();
    const password = document.getElementById('reg_password').value;
    const sex = parseInt(document.getElementById('reg_sex').value);
    const grade = document.getElementById('reg_grade').value;
    const age = parseInt(document.getElementById('reg_age').value);
    
    // 清除之前的错误信息
    document.querySelectorAll('.error-message').forEach(el => el.textContent = '');
    
    // 数据校验
    let hasError = false;
    
    const nameError = validateName(name);
    if (nameError) {
        document.getElementById('reg_name_error').textContent = nameError;
        hasError = true;
    }
    
    const stuIdError = validateStudentID(stuId);
    if (stuIdError) {
        document.getElementById('reg_stu_id_error').textContent = stuIdError;
        hasError = true;
    }
    
    const classError = validateClass(stuClass);
    if (classError) {
        document.getElementById('reg_stu_class_error').textContent = classError;
        hasError = true;
    }
    
    const passwordError = validatePassword(password);
    if (passwordError) {
        document.getElementById('reg_password_error').textContent = passwordError;
        hasError = true;
    }
    
    const ageError = validateAge(age);
    if (ageError) {
        document.getElementById('reg_age_error').textContent = ageError;
        hasError = true;
    }
    
    if (hasError) {
        return;
    }
    
    const registerBtn = document.querySelector('.register-container .login-btn');
    showLoading(registerBtn);
    
    try {
        const response = await apiRequest('/v1/api/public/register', {
            method: 'POST',
            body: JSON.stringify({
                name: name,
                student_id: stuId,
                student_class: stuClass,
                password: password,
                sex: sex,
                grade: grade,
                age: age
            })
        });
        
        alert('注册成功！请返回登录。');
        showLogin();
        
    } catch (error) {
        try {
            const errorData = JSON.parse(error.message);
            showErrorModal(errorData.status, errorData.info);
        } catch {
            showErrorModal('错误', error.message);
        }
    } finally {
        hideLoading(registerBtn);
    }
}

function showRegister() {
    document.querySelector('.login-container').style.display = 'none';
    document.querySelector('.register-container').style.display = 'block';
}

function showLogin() {
    document.querySelector('.register-container').style.display = 'none';
    document.querySelector('.login-container').style.display = 'block';
}

async function logout() {
    try {
        await apiRequest('/v1/api/stu-manager/stu-logout', {
            method: 'GET'
        });
    } catch (error) {
        console.error('Logout error:', error);
    }
    
    localStorage.removeItem('authToken');
    localStorage.removeItem('userRole');
    window.location.href = 'index.html';
}

// 检查用户认证状态，防止无限刷新
function checkAuth() {
    // 防止重复检查
    if (window.isCheckingAuth) return;
    window.isCheckingAuth = true;
    
    try {
        authToken = localStorage.getItem('authToken');
        const userRole = localStorage.getItem('userRole');
        const currentPage = window.location.pathname.split('/').pop();
        
        // 未登录状态
        if (!authToken) {
            if (currentPage !== 'index.html' && currentPage !== '') {
                window.location.href = 'index.html';
            }
            return;
        }
        
        // 已登录但页面不匹配 - 只在确实需要跳转时才跳转
        if (userRole === 'admin' && currentPage !== 'admin.html' && currentPage !== '') {
            window.location.href = 'admin.html';
        } else if (userRole === 'student' && currentPage !== 'student.html' && currentPage !== '') {
            window.location.href = 'student.html';
        }
    } finally {
        // 确保检查完成后重置标志
        setTimeout(() => {
            window.isCheckingAuth = false;
        }, 100);
    }
}




// 学生页面功能
function showSection(sectionId) {
    // 隐藏所有section
    document.querySelectorAll('.section').forEach(section => {
        section.classList.remove('active');
    });
    
    // 显示选中的section
    document.getElementById(sectionId).classList.add('active');
    
    // 更新菜单状态
    document.querySelectorAll('.menu-item').forEach(item => {
        item.classList.remove('active');
    });
    
    event.target.closest('.menu-item').classList.add('active');
    currentSection = sectionId;
    
    // 根据section加载数据
    if (sectionId === 'selectable-courses') {
        loadSelectableCourses();
    } else if (sectionId === 'selected-courses') {
        loadSelectedCourses();
    } else if (sectionId === 'profile') {
        loadStudentInfo();
    }
}

async function loadStudentInfo() {
    try {
        const response = await apiRequest('/v1/api/stu-manager/stu-info', {
            method: 'GET'
        });
        
        const student = response.data;
        currentUser = student;
        
        // 更新导航栏显示的学生姓名
        const studentNameElement = document.getElementById('studentName');
        if (studentNameElement) {
            studentNameElement.textContent = student.name;
        }
        
        // 填充个人信息表单
        document.getElementById('profile_name').value = student.name;
        document.getElementById('profile_stu_id').value = student.stu_id || student.student_id;
        document.getElementById('profile_stu_class').value = student.stu_class || student.student_class;
        document.getElementById('profile_grade').value = student.grade;
        document.getElementById('profile_sex').value = student.sex;
        document.getElementById('profile_age').value = student.age;
        
    } catch (error) {
        try {
            const errorData = JSON.parse(error.message);
            showErrorModal(errorData.status, errorData.info);
        } catch {
            showErrorModal('错误', error.message);
        }
    }
}

async function loadSelectableCourses() {
    try {
        const response = await apiRequest('/v1/api/class-manager/get-selectable-classes', {
            method: 'GET'
        });
        
        const courses = response.data.selectable_classes || [];
        const coursesList = document.getElementById('selectableCoursesList');
        
        if (courses.length === 0) {
            coursesList.innerHTML = `
                <div class="empty-state">
                    <h3>暂无可选课程</h3>
                    <p>当前没有可选的课程，请稍后再试。</p>
                </div>
            `;
            return;
        }
        
        coursesList.innerHTML = courses.map(course => `
            <div class="course-card">
                <h3>${course.class_name}</h3>
                <div class="course-info"><strong>课程ID：</strong>${course.class_id}</div>
                <div class="course-info"><strong>授课教师：</strong>${course.class_teacher}</div>
                <div class="course-info"><strong>上课地点：</strong>${course.class_location}</div>
                <div class="course-info"><strong>上课时间：</strong>${course.class_time}</div>
                <div class="course-info"><strong>课程容量：</strong>${course.class_selsetion}/${course.class_capcity}</div>
                <div class="course-actions">
                    <button class="select-btn" onclick="subscribeCourse('${course.class_id}')" 
                            ${course.class_selsetion >= course.class_capcity ? 'disabled' : ''}>
                        ${course.class_selsetion >= course.class_capcity ? '已满' : '选课'}
                    </button>
                </div>
            </div>
        `).join('');
        
    } catch (error) {
        try {
            const errorData = JSON.parse(error.message);
            showErrorModal(errorData.status, errorData.info);
        } catch {
            showErrorModal('错误', error.message);
        }
    }
}

async function loadSelectedCourses() {
    try {
        const response = await apiRequest('/v1/api/class-manager/get-subscribed-classes', {
            method: 'GET'
        });
        
        const courses = response.data.selected_classes || [];
        const coursesList = document.getElementById('selectedCoursesList');
        
        if (courses.length === 0) {
            coursesList.innerHTML = `
                <div class="empty-state">
                    <h3>暂无已选课程</h3>
                    <p>您还没有选择任何课程，请前往可选课程页面选课。</p>
                </div>
            `;
            return;
        }
        
        coursesList.innerHTML = courses.map(course => `
            <div class="course-card">
                <h3>${course.class_name}</h3>
                <div class="course-info"><strong>课程ID：</strong>${course.class_id}</div>
                <div class="course-info"><strong>授课教师：</strong>${course.class_teacher}</div>
                <div class="course-info"><strong>上课地点：</strong>${course.class_location}</div>
                <div class="course-info"><strong>上课时间：</strong>${course.class_time}</div>
                <div class="course-actions">
                    <button class="drop-btn" onclick="dropCourse('${course.class_id}')">退课</button>
                </div>
            </div>
        `).join('');
        
    } catch (error) {
        try {
            const errorData = JSON.parse(error.message);
            showErrorModal(errorData.status, errorData.info);
        } catch {
            showErrorModal('错误', error.message);
        }
    }
}

async function subscribeCourse(classId) {
    showConfirmModal('确定要选择这门课程吗？', async () => {
        try {
            const response = await apiRequest('/v1/api/class-manager/subscribe-class', {
                method: 'POST',
                body: JSON.stringify({
                    class_id: classId
                })
            });
            
            alert('选课成功！');
            loadSelectableCourses();
            
        } catch (error) {
            try {
                const errorData = JSON.parse(error.message);
                showErrorModal(errorData.status, errorData.info);
            } catch {
                showErrorModal('错误', error.message);
            }
        }
    });
}

async function dropCourse(classId) {
    showConfirmModal('确定要退选这门课程吗？', async () => {
        try {
            const response = await apiRequest('/v1/api/class-manager/del-class', {
                method: 'DELETE',
                body: JSON.stringify({
                    class_id: classId
                })
            });
            
            alert('退课成功！');
            loadSelectedCourses();
            
        } catch (error) {
            try {
                const errorData = JSON.parse(error.message);
                showErrorModal(errorData.status, errorData.info);
            } catch {
                showErrorModal('错误', error.message);
            }
        }
    });
}

function enableEdit() {
    const inputs = document.querySelectorAll('#profile input, #profile select');
    inputs.forEach(input => {
        input.disabled = false;
    });
    document.getElementById('profileActions').style.display = 'flex';
}

function cancelEdit() {
    const inputs = document.querySelectorAll('#profile input, #profile select');
    inputs.forEach(input => {
        input.disabled = true;
    });
    document.getElementById('profileActions').style.display = 'none';
    loadStudentInfo(); // 重新加载数据
}

async function saveProfile() {
    // 清除之前的错误信息
    document.querySelectorAll('#profile .error-message').forEach(el => el.textContent = '');
    
    // 获取表单数据
    const name = document.getElementById('profile_name').value.trim();
    const stuClass = document.getElementById('profile_stu_class').value.trim();
    const grade = document.getElementById('profile_grade').value;
    const sex = parseInt(document.getElementById('profile_sex').value);
    const age = parseInt(document.getElementById('profile_age').value);
    
    // 数据校验
    let hasError = false;
    
    const nameError = validateName(name);
    if (nameError) {
        document.getElementById('profile_name_error').textContent = nameError;
        hasError = true;
    }
    
    const classError = validateClass(stuClass);
    if (classError) {
        document.getElementById('profile_stu_class_error').textContent = classError;
        hasError = true;
    }
    
    const ageError = validateAge(age);
    if (ageError) {
        document.getElementById('profile_age_error').textContent = ageError;
        hasError = true;
    }
    
    if (hasError) {
        return;
    }
    
    // 构建更新字段
    const updateColumns = [];
    
    if (name !== currentUser.name) {
        updateColumns.push({ field: 'name', value: name });
    }
    if (stuClass !== (currentUser.stu_class || currentUser.student_class)) {
        updateColumns.push({ field: 'student_class', value: stuClass });
    }
    if (grade !== currentUser.grade) {
        updateColumns.push({ field: 'grade', value: grade });
    }
    if (sex !== currentUser.sex) {
        updateColumns.push({ field: 'sex', value: sex.toString() });
    }
    if (age !== currentUser.age) {
        updateColumns.push({ field: 'age', value: age.toString() });
    }
    
    if (updateColumns.length === 0) {
        alert('没有需要更新的信息');
        cancelEdit();
        return;
    }
    
    try {
        const response = await apiRequest('/v1/api/stu-manager/stu-update', {
            method: 'PATCH',
            body: JSON.stringify({
                update_columns: updateColumns
            })
        });
        
        alert('个人信息更新成功！');
        cancelEdit();
        
    } catch (error) {
        try {
            const errorData = JSON.parse(error.message);
            showErrorModal(errorData.status, errorData.info);
        } catch {
            showErrorModal('错误', error.message);
        }
    }
}

// 管理员页面功能
async function loadStudentList(page = 1, pageSize = 15) {
    try {
        const response = await apiRequest(`/v1/api/admin/stu-manager/get-stu-list?page=${page}&resNum=${pageSize}`, {
            method: 'GET'
        });
        
        const { students_list, total } = response.data;
        const tbody = document.getElementById('studentTableBody');
        
        if (students_list.length === 0) {
            tbody.innerHTML = `
                <tr>
                    <td colspan="5" style="text-align: center; padding: 40px; color: #666;">
                        暂无学生数据
                    </td>
                </tr>
            `;
            return;
        }
        
        tbody.innerHTML = students_list.map(student => `
            <tr>
                <td>${student.stu_id}</td>
                <td>${student.student_name || student.name}</td>
                <td>${student.student_class || student.stu_class}</td>
                <td>${student.grade}</td>
                <td>
                    <div class="action-buttons">
                        <button class="view-btn" onclick="viewStudentCourses('${student.stu_id}')">查看选课</button>
                        <button class="edit-btn-small" onclick="editStudent('${student.stu_id}')">编辑</button>
                        <button class="delete-btn" onclick="deleteStudent('${student.stu_id}')">删除</button>
                    </div>
                </td>
            </tr>
        `).join('');
        
        // 生成分页
        generatePagination('studentPagination', total, page, pageSize, loadStudentList);
        
    } catch (error) {
        try {
            const errorData = JSON.parse(error.message);
            showErrorModal(errorData.status, errorData.info);
        } catch {
            showErrorModal('错误', error.message);
        }
    }
}

async function loadAdminCourses() {
    try {
        const response = await apiRequest('/v1/api/admin/classes-manager/get-class-status', {
            method: 'GET'
        });
        
        const courses = response.data.selectable_classes || [];
        const tbody = document.getElementById('courseTableBody');
        
        if (courses.length === 0) {
            tbody.innerHTML = `
                <tr>
                    <td colspan="8" style="text-align: center; padding: 40px; color: #666;">
                        暂无课程数据
                    </td>
                </tr>
            `;
            return;
        }
        
        tbody.innerHTML = courses.map(course => `
            <tr>
                <td>${course.class_id}</td>
                <td>${course.class_name}</td>
                <td>${course.class_teacher}</td>
                <td>${course.class_location}</td>
                <td>${course.class_time}</td>
                <td>${course.class_capcity}</td>
                <td>${course.class_selsetion}</td>
                <td>
                    <div class="action-buttons">
                        <button class="edit-btn-small" onclick="editCourse('${course.class_id}')">编辑</button>
                        <button class="delete-btn" onclick="deleteCourse('${course.class_id}')">删除</button>
                    </div>
                </td>
            </tr>
        `).join('');
        
    } catch (error) {
        try {
            const errorData = JSON.parse(error.message);
            showErrorModal(errorData.status, errorData.info);
        } catch {
            showErrorModal('错误', error.message);
        }
    }
}

function generatePagination(containerId, total, currentPage, pageSize, loadFunction) {
    const totalPages = Math.ceil(total / pageSize);
    const container = document.getElementById(containerId);
    
    if (totalPages <= 1) {
        container.innerHTML = '';
        return;
    }
    
    let paginationHTML = '';
    
    // 上一页
    if (currentPage > 1) {
        paginationHTML += `<button onclick="${loadFunction.name}(${currentPage - 1}, ${pageSize})">上一页</button>`;
    }
    
    // 页码
    for (let i = 1; i <= totalPages; i++) {
        if (i === currentPage) {
            paginationHTML += `<button class="active" disabled>${i}</button>`;
        } else if (Math.abs(i - currentPage) <= 2 || i === 1 || i === totalPages) {
            paginationHTML += `<button onclick="${loadFunction.name}(${i}, ${pageSize})">${i}</button>`;
        } else if (Math.abs(i - currentPage) === 3) {
            paginationHTML += `<span>...</span>`;
        }
    }
    
    // 下一页
    if (currentPage < totalPages) {
        paginationHTML += `<button onclick="${loadFunction.name}(${currentPage + 1}, ${pageSize})">下一页</button>`;
    }
    
    container.innerHTML = paginationHTML;
}

function showCreateStudent() {
    document.getElementById('studentEditTitle').textContent = '创建学生';
    document.getElementById('edit_stu_id').readOnly = false;
    
    // 清空表单
    document.getElementById('edit_stu_name').value = '';
    document.getElementById('edit_stu_id').value = '';
    document.getElementById('edit_stu_class').value = '';
    document.getElementById('edit_grade').value = '大一';
    document.getElementById('edit_sex').value = '1';
    document.getElementById('edit_age').value = '';
    document.getElementById('edit_password').value = '';
    
    document.getElementById('studentEditModal').style.display = 'block';
}

function editStudent(stuId) {
    document.getElementById('studentEditTitle').textContent = '编辑学生信息';
    document.getElementById('edit_stu_id').readOnly = true;
    document.getElementById('edit_stu_id').value = stuId;
    
    // 这里应该调用API获取学生详细信息
    // 为了简化，我们假设数据已经在表格中
    const row = event.target.closest('tr');
    document.getElementById('edit_stu_name').value = row.cells[1].textContent;
    document.getElementById('edit_stu_class').value = row.cells[2].textContent;
    document.getElementById('edit_grade').value = row.cells[3].textContent;
    
    document.getElementById('studentEditModal').style.display = 'block';
}

async function saveStudentEdit() {
    const isCreate = document.getElementById('studentEditTitle').textContent === '创建学生';
    const stuId = document.getElementById('edit_stu_id').value.trim();
    const name = document.getElementById('edit_stu_name').value.trim();
    const stuClass = document.getElementById('edit_stu_class').value.trim();
    const grade = document.getElementById('edit_grade').value;
    const sex = parseInt(document.getElementById('edit_sex').value);
    const age = parseInt(document.getElementById('edit_age').value);
    const password = document.getElementById('edit_password').value;
    
    // 数据校验
    let hasError = false;
    document.querySelectorAll('[id$="_error"]').forEach(el => el.textContent = '');
    
    if (isCreate) {
        const stuIdError = validateStudentID(stuId);
        if (stuIdError) {
            // 创建时没有专门的错误显示元素，使用alert
            alert(stuIdError);
            hasError = true;
        }
    }
    
    const nameError = validateName(name);
    if (nameError) {
        document.getElementById('edit_stu_name_error').textContent = nameError;
        hasError = true;
    }
    
    const classError = validateClass(stuClass);
    if (classError) {
        document.getElementById('edit_stu_class_error').textContent = classError;
        hasError = true;
    }
    
    if (isCreate && !password) {
        document.getElementById('edit_password_error').textContent = '密码不能为空';
        hasError = true;
    } else if (password) {
        const passwordError = validatePassword(password);
        if (passwordError) {
            document.getElementById('edit_password_error').textContent = passwordError;
            hasError = true;
        }
    }
    
    const ageError = validateAge(age);
    if (ageError) {
        document.getElementById('edit_age_error').textContent = ageError;
        hasError = true;
    }
    
    if (hasError) {
        return;
    }
    
    try {
        if (isCreate) {
            await apiRequest('/v1/api/admin/stu-manager/create-stu', {
                method: 'POST',
                body: JSON.stringify({
                    name: name,
                    student_id: stuId,
                    student_class: stuClass,
                    password: password,
                    sex: sex,
                    grade: grade,
                    age: age
                })
            });
            alert('学生创建成功！');
        } else {
            const updateColumns = [];
            updateColumns.push({ field: 'name', value: name });
            updateColumns.push({ field: 'student_class', value: stuClass });
            updateColumns.push({ field: 'grade', value: grade });
            updateColumns.push({ field: 'sex', value: sex.toString() });
            updateColumns.push({ field: 'age', value: age.toString() });
            
            if (password) {
                updateColumns.push({ field: 'password', value: password });
            }
            
            await apiRequest('/v1/api/admin/stu-manager/update-stu-info', {
                method: 'PATCH',
                body: JSON.stringify({
                    stu_id: stuId,
                    update_columns: updateColumns
                })
            });
            alert('学生信息更新成功！');
        }
        
        closeStudentEditModal();
        loadStudentList();
        
    } catch (error) {
        try {
            const errorData = JSON.parse(error.message);
            showErrorModal(errorData.status, errorData.info);
        } catch {
            showErrorModal('错误', error.message);
        }
    }
}

function closeStudentEditModal() {
    document.getElementById('studentEditModal').style.display = 'none';
}

async function deleteStudent(stuId) {
    showConfirmModal('确定要删除这个学生吗？此操作不可恢复。', async () => {
        try {
            await apiRequest('/v1/api/admin/stu-manager/del-stu', {
                method: 'DELETE',
                body: JSON.stringify({
                    stu_id: stuId
                })
            });
            
            alert('学生删除成功！');
            loadStudentList();
            
        } catch (error) {
            try {
                const errorData = JSON.parse(error.message);
                showErrorModal(errorData.status, errorData.info);
            } catch {
                showErrorModal('错误', error.message);
            }
        }
    });
}

function showCreateCourse() {
    document.getElementById('courseEditTitle').textContent = '创建课程';
    document.getElementById('edit_class_id').readOnly = false;
    
    // 清空表单
    document.getElementById('edit_class_name').value = '';
    document.getElementById('edit_class_id').value = '';
    document.getElementById('edit_class_teacher').value = '';
    document.getElementById('edit_class_location').value = '';
    document.getElementById('edit_class_time').value = '';
    document.getElementById('edit_class_capacity').value = '';
    
    document.getElementById('courseEditModal').style.display = 'block';
}

function editCourse(courseId) {
    document.getElementById('courseEditTitle').textContent = '编辑课程信息';
    document.getElementById('edit_class_id').readOnly = true;
    document.getElementById('edit_class_id').value = courseId;
    
    // 这里应该调用API获取课程详细信息
    // 为了简化，我们假设数据已经在表格中
    const row = event.target.closest('tr');
    document.getElementById('edit_class_name').value = row.cells[1].textContent;
    document.getElementById('edit_class_teacher').value = row.cells[2].textContent;
    document.getElementById('edit_class_location').value = row.cells[3].textContent;
    document.getElementById('edit_class_time').value = row.cells[4].textContent;
    document.getElementById('edit_class_capacity').value = row.cells[5].textContent;
    
    document.getElementById('courseEditModal').style.display = 'block';
}

async function saveCourseEdit() {
    const isCreate = document.getElementById('courseEditTitle').textContent === '创建课程';
    const classId = document.getElementById('edit_class_id').value.trim();
    const className = document.getElementById('edit_class_name').value.trim();
    const teacher = document.getElementById('edit_class_teacher').value.trim();
    const location = document.getElementById('edit_class_location').value.trim();
    const time = document.getElementById('edit_class_time').value.trim();
    const capacity = parseInt(document.getElementById('edit_class_capacity').value);
    
    // 简单校验
    if (!classId || !className || !teacher || !location || !time || !capacity) {
        alert('请填写所有必填字段');
        return;
    }
    
    try {
        if (isCreate) {
            await apiRequest('/v1/api/admin/classes-manager/add-course', {
                method: 'POST',
                body: JSON.stringify({
                    class_name: className,
                    class_id: classId,
                    class_location: location,
                    class_time: time,
                    class_teacher: teacher,
                    class_capcity: capacity
                })
            });
            alert('课程创建成功！');
        } else {
            const updateColumns = [
                { field: 'class_name', value: className },
                { field: 'class_teacher', value: teacher },
                { field: 'class_location', value: location },
                { field: 'class_time', value: time }
            ];
            
            await apiRequest('/v1/api/admin/classes-manager/edit-class-info', {
                method: 'PATCH',
                body: JSON.stringify({
                    class_id: classId,
                    update_columns: updateColumns
                })
            });
            
            // 更新容量
            await apiRequest('/v1/api/admin/classes-manager/edit-class-stock', {
                method: 'PATCH',
                body: JSON.stringify({
                    class_id: classId,
                    stock: capacity
                })
            });
            
            alert('课程信息更新成功！');
        }
        
        closeCourseEditModal();
        loadAdminCourses();
        
    } catch (error) {
        try {
            const errorData = JSON.parse(error.message);
            showErrorModal(errorData.status, errorData.info);
        } catch {
            showErrorModal('错误', error.message);
        }
    }
}

function closeCourseEditModal() {
    document.getElementById('courseEditModal').style.display = 'none';
}

async function deleteCourse(courseId) {
    showConfirmModal('确定要删除这门课程吗？', async () => {
        try {
            await apiRequest('/v1/api/admin/classes-manager/delete-course', {
                method: 'DELETE',
                body: JSON.stringify({
                    class_id: courseId
                })
            });
            
            alert('课程删除成功！');
            loadAdminCourses();
            
        } catch (error) {
            try {
                const errorData = JSON.parse(error.message);
                showErrorModal(errorData.status, errorData.info);
            } catch {
                showErrorModal('错误', error.message);
            }
        }
    });
}

async function viewStudentCourses(stuId) {
    try {
        const response = await apiRequest(`/v1/api/admin/classes-manager/get-stu-classes/${stuId}`, {
            method: 'GET'
        });
        
        const courses = response.data.selected_classes || [];
        let message = `学生 ${stuId} 的选课情况：\n\n`;
        
        if (courses.length === 0) {
            message += '该学生暂未选择任何课程。';
        } else {
            courses.forEach(course => {
                message += `• ${course.class_name} (${course.class_id})\n`;
                message += `  教师：${course.class_teacher}\n`;
                message += `  时间：${course.class_time}\n\n`;
            });
        }
        
        alert(message);
        
    } catch (error) {
        try {
            const errorData = JSON.parse(error.message);
            showErrorModal(errorData.status, errorData.info);
        } catch {
            showErrorModal('错误', error.message);
        }
    }
}

async function startSelection() {
    try {
        await apiRequest('/v1/api/admin/classes-manager/start-course-select-event', {
            method: 'GET'
        });
        
        alert('选课已开始！');
        document.getElementById('selectionStatus').textContent = '进行中';
        document.getElementById('selectionStatus').className = 'status-running';
        
    } catch (error) {
        try {
            const errorData = JSON.parse(error.message);
            showErrorModal(errorData.status, errorData.info);
        } catch {
            showErrorModal('错误', error.message);
        }
    }
}

async function stopSelection() {
    try {
        await apiRequest('/v1/api/admin/classes-manager/stop-course-select-event', {
            method: 'GET'
        });
        
        alert('选课已停止！');
        document.getElementById('selectionStatus').textContent = '已停止';
        document.getElementById('selectionStatus').className = 'status-stopped';
        
    } catch (error) {
        try {
            const errorData = JSON.parse(error.message);
            showErrorModal(errorData.status, errorData.info);
        } catch {
            showErrorModal('错误', error.message);
        }
    }
}

// 页面加载完成后的初始化 - 只处理主页功能，避免重复初始化
document.addEventListener('DOMContentLoaded', function() {
    // 根据当前页面进行相应初始化
    const currentPage = window.location.pathname.split('/').pop();
    
    // 只在主页设置模态框和键盘事件
    if (currentPage === 'index.html' || currentPage === '') {
        // 为模态框添加点击外部关闭功能
        document.addEventListener('click', function(event) {
            if (event.target.classList.contains('modal')) {
                event.target.style.display = 'none';
            }
        });
        
        // 为回车键添加登录功能
        document.addEventListener('keypress', function(event) {
            if (event.key === 'Enter') {
                const loginBtn = document.querySelector('.login-btn');
                if (loginBtn && !loginBtn.disabled) {
                    login();
                }
            }
        });
    }
});