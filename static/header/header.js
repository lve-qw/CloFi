const dropdownMenu = document.querySelector('.dropdown-menu');

function logoutUser() {
    localStorage.removeItem('token');
    window.location.reload();
}

function updateMenu() {
    const token = localStorage.getItem('token');
    
    if (token) {
        try {
            const tokenParts = token.split('.');
            if (tokenParts.length !== 3) {
                throw new Error('Неверный формат токена')
            }
            const tokenData = JSON.parse(atob(tokenParts[1]));
            const userName = tokenData.name || 'Пользователь';
            
            dropdownMenu.innerHTML = `
                <li><a class="dropdown-item" href="#">Привет, ${userName}</a></li>
                <li><a class="dropdown-item" href="#" id="logoutBtn">Выйти</a></li>
            `;
            
            document.getElementById('logoutBtn').addEventListener('click', function(e) {
                e.preventDefault();
                logoutUser();
            });
        } catch (error) {
            localStorage.removeItem('token');
            updateMenu();
        }
    } else {
        dropdownMenu.innerHTML = `
            <li><a class="dropdown-item" href="../login/login.html">Войти</a></li>
            <li><a class="dropdown-item" href="../register/register.html">Регистрация</a></li>
        `;
    }
}
updateMenu();