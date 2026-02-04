const loginForm = document.getElementById('loginForm');

function showError(message) {
    const errorContainer = loginForm.querySelector('.login-error');
    errorContainer.textContent = message;
    errorContainer.style.display = 'block';
}

function validation(username, password) {
    if (!username || !password) {
        return 'Не все поля заполнены, попробуйте еще раз';
    }
    if (username.length < 3 || username.length > 32) {
        return 'Логин должен состоять из 3-32 символов';
    }
    if ((/[^a-z0-9]/i).test(username)) {
        return 'Логин должен состоять только из латинских букв и цифр';
    }
    if (password.length < 6) {
        return 'Пароль должен быть минимум из 6 символов';
    }
    if (password.length > 128) {
        return 'Пароль слишком длинный (максимум 128 символов)';
    }
    if (password.includes(' ')) {
        return 'Пароль не должен содержать пробелов';
    }
    return false;
}

loginForm.addEventListener('submit', async (event) => {
    event.preventDefault();
    const errorContainer = loginForm.querySelector('.login-error');
    errorContainer.textContent = '';
    errorContainer.style.display = 'none';

    const formData = {
        username: document.getElementById('username').value.trim(),
        password: document.getElementById('password').value.trim()
    };

    const validationError = validation(formData.username, formData.password);
    if (validationError) {
        showError(validationError);
        return;
    }

    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(formData)
        });

        const data = await response.json();

        if (response.ok) {
            localStorage.setItem('token', data.token);
            window.location.href = '../main/main.html';
        } else if (response.status === 401) {
            showError('Неверный логин или пароль');
        } else {
            showError(data.error || 'Произошла ошибка');
        }
    } catch (error) {
        console.error('Ошибка соединения:', error);
        showError('Ошибка соединения с сервером');
    }
});