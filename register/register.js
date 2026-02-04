const registerForm = document.getElementById('registerForm');

function showError(message) {
    const errorContainer = registerForm.querySelector(`.register-error`);
    errorContainer.textContent = message; 
    errorContainer.style.display = 'block';
}

function validation(username, name, password) {
    if (!username || !name || !password) {
        return 'Не все поля заполнены, попробуйте еще раз';
    }

    if (username.length < 3 || username.length > 32) {
        return 'Логин должен состоять из 3-32 символов';
    } else if ((/[^a-z0-9]/i).test(username)) {
        return 'Логин должен состоять только из латинских букв и цифр';
    } else if (name.length < 1 || name.length > 64) {
        return 'Имя должно состоять из 1-64 символов';
    } else if (password.length < 6) {
        return 'Пароль должен быть минимум из 6 символов';
    } else if (password.length > 128) {
        return 'Пароль слишком длинный (максимум 128 символов)';
    } else if (password.includes(' ')) {
        return 'Пароль не должен содержать пробелов';
    }
    return false;
}

registerForm.addEventListener(`submit`, async (event) => {
    event.preventDefault();
    const errorContainer = registerForm.querySelector(`.register-error`);
    errorContainer.textContent = '';  
    errorContainer.style.display = 'none';

    const formData = {
        username: document.getElementById('username').value.trim(),
        name: document.getElementById('name').value.trim(),
        password: document.getElementById('password').value.trim()
    };

    const validationError = validation(formData.username, formData.name, formData.password);
    if (validationError) {
        showError(validationError);
        return;
    }

    try {
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(formData)
        });

        const data = await response.json();

        if (response.ok) {
            if (response.status === 201) {
                window.location.href = '../login/login.html';
            } else {
                showError(data.error || 'Произошла ошибка');
            }
        } else {
            if (response.status === 409) {
                showError('Пользователь с таким логином уже существует');
            } else if (response.status === 400) {
                showError(data.error || 'Неверный формат данных');
            } else {
                console.error('Register error:', response.status, data);
                showError(data.error || 'Произошла ошибка');
            }
        }

    } catch (error) {
        console.error('Ошибка соединения с сервером:', error);
        showError('Ошибка соединения с сервером');
    }
});

