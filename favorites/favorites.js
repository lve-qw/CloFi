const favoritesContainer = document.getElementById('favorites');
const messageContainer = document.querySelector('.message');

async function loadFavorites() {
  const token = localStorage.getItem('token');

  if (!token) {
    if (favoritesContainer) {
      favoritesContainer.innerHTML = '';
    }
    return;
  }

  try {
    const res = await fetch('/api/likes', {
      headers: { 'Authorization': 'Bearer ' + token }
    });

    if (!res.ok) {
      if (res.status === 401) {
        localStorage.removeItem('token');
      }
      const err = await res.json().catch(() => ({ error: 'Не удалось загрузить избранные товары' }));
      messageContainer.textContent = err.error || 'Не удалось загрузить избранные товары';
      favoritesContainer.innerHTML = '';
      return;
    }

    const products = await res.json();

    if (!products || products.length === 0) {
      messageContainer.textContent = 'У вас пока нет избранных товаров';
      favoritesContainer.innerHTML = '';
      return;
    }

    messageContainer.textContent = '';
    renderFavorites(products);

  } catch (error) {
    console.error('Ошибка:', error);
    messageContainer.textContent = 'Не удалось загрузить избранные товары';
    favoritesContainer.innerHTML = '';
  }
}

function renderFavorites(products) {
  if (!favoritesContainer) return;
  favoritesContainer.innerHTML = '';

  products.forEach(product => {
    const photo = product.photos_urls && product.photos_urls[0] ? product.photos_urls[0] : '../assets/placeholder.JPG';

    const productCard = document.createElement('div');
    productCard.className = 'product-card';
    productCard.dataset.id = product.id;

    productCard.innerHTML = `
            <img src="${photo}" alt="${product.name}" class="product-image">
        <div class="favorite-icon" data-id="${product.id}">
            <svg width="38" height="38" viewBox="0 0 38 38" fill="none">
            <path d="M32.9967 7.29917C32.188 6.4901 31.2278 5.84828 30.171 5.41039C29.1142 4.9725 27.9814 4.74712 26.8375 4.74712C25.6936 4.74712 24.5608 4.9725 23.504 5.41039C22.4472 5.84828 21.487 6.4901 20.6783 7.29917L19 8.97751L17.3217 7.29917C15.6882 5.66566 13.4726 4.74796 11.1625 4.74796C8.85237 4.74796 6.63685 5.66566 5.00334 7.29917C3.36982 8.93269 2.45213 11.1482 2.45213 13.4583C2.45213 15.7685 3.36982 17.984 5.00334 19.6175L19 33.6142L32.9967 19.6175C33.8057 18.8088 34.4476 17.8486 34.8854 16.7918C35.3233 15.735 35.5487 14.6023 35.5487 13.4583C35.5487 12.3144 35.3233 11.1817 34.8854 10.1249C34.4476 9.06805 33.8057 8.10787 32.9967 7.29917Z" stroke="${isLiked ? '#5A87E6': '#5A5A6E'}" fill="${isLiked ? '#5A87E6' : 'none'}" stroke-width="4"/>
            </svg>
        </div>
        </div>
        <div class="card-body">
        <div class="product-title-brand">
            <span class="product-name">${product.name}</span>
            ${product.brand ? `<span class="product-brand">${product.brand}</span>` : ''}
        </div>
        <span class="price">${product.price} ₽</span>
        </div>
        `;

    favoritesContainer.appendChild(productCard);
    productCard.addEventListener('click', () => {
        const modal = new bootstrap.Modal(document.getElementById('productModal'));
        document.getElementById('modalAvailability').textContent = product.availability ? 'В наличии' : 'Нет в наличии';
        document.getElementById('modalDescription').textContent = product.description || '-';
        document.getElementById('modalLink').href = product.url || '#';
        modal.show();
    });
  });

  addFavoriteHandlers();
}

if (favoritesContainer) {
  favoritesContainer.addEventListener('click', async (event) => {
    const icon = event.target.closest('.favorite-icon');
    if (!icon) return; 

    event.stopPropagation(); 

    const productId = icon.dataset.id;
    const token = localStorage.getItem('token');

    if (!token) {
      messageContainer.textContent = 'Войдите в аккаунт, чтобы управлять избранным';
      return;
    }

    try {
      const res = await fetch(`/api/like?product_id=${productId}`, {
        method: 'POST',
        headers: { 'Authorization': 'Bearer ' + token }
      });

      if (!res.ok) {
        if (res.status === 401) localStorage.removeItem('token');
        const err = await res.json().catch(() => ({ error: 'Ошибка сервера' }));
        messageContainer.textContent = err.error || 'Ошибка при обновлении избранного';
        return;
      }

      const data = await res.json();
      const heart = icon.querySelector('path');

      if (data.status === 'лайк добавлен') {
        heart.setAttribute('fill', '#5A87E6');
        heart.setAttribute('stroke', '#5A87E6');
      } else if (data.status === 'лайк удалён') {
        heart.removeAttribute('fill');
        heart.setAttribute('stroke', '#5A5A6E');
      }
    } catch (error) {
      messageContainer.textContent = 'Ошибка при обновлении избранного';
    }
  });
}

loadFavorites();
