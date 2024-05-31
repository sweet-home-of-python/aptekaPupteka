document.addEventListener('DOMContentLoaded', () => {
    const medicationsList = document.getElementById('medicationsList');
    const searchInput = document.getElementById('search');
    const addBtn = document.getElementById('addBtn');
    const modal = document.getElementById('modal');
    const adjustModal = document.getElementById('adjustModal');
    const closeBtns = document.querySelectorAll('.close');
    const medicationForm = document.getElementById('medicationForm');
    const adjustForm = document.getElementById('adjustForm');
    const prevPage = document.getElementById('prevPage');
    const nextPage = document.getElementById('nextPage');
    const pageIndicator = document.getElementById('pageIndicator');

    let currentPage = 0;
    const itemsPerPage = 5;
    let medications = [];
    let adjustMedicationId = null;
    let adjustAction = '';

    function fetchMedications(page = currentPage, limit = itemsPerPage) {
    const url = '/getPage';

    const requestData = {
        page: page,
        limit: limit
    };

    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
    })
    .then(response => response.json())
    .then(data => {
        medications = data;
        renderMedications();
    })
    .catch(error => console.error('Error fetching medications:', error));
}

    function renderMedications() {
        medicationsList.innerHTML = '';
        medications.forEach(medication => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${medication.name}</td>
                <td>${medication.quantity}</td>
                

                <td>
                  
                    <button onclick="deleteMedication(\'${medication.name}\')">Удалить</button>
                    <button onclick="showAdjustModal(\'${medication.name}\', 'add')">Добавить</button>

                </td>
                
            `;
            medicationsList.appendChild(row);
        });
        pageIndicator.textContent = `Страница ${currentPage}`;
    }

    searchInput.addEventListener('input', () => {
        fetchMedications(searchInput.value, currentPage);
    });

    addBtn.addEventListener('click', () => {
        modal.style.display = 'block';
        medicationForm.reset();
    });

    closeBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            modal.style.display = 'none';
            adjustModal.style.display = 'none';
        });
    });

    medicationForm.addEventListener('submit', (event) => {
        event.preventDefault();
        console.log('Form submitted');  // Логирование события

        const formData = new FormData(medicationForm);
        const medication = {
            drug: formData.get('name'),
        };

        fetch('http://localhost:8082/newDrug', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(medication)
        })
        .then(response => {
            console.log(response);  // Логирование ответа
            return response.json();
        })
        .then(() => {
            modal.style.display = 'none';
            fetchMedications(searchInput.value, currentPage);
        })
        .catch(error => {
            console.error('Error:', error);  // Логирование ошибок
        });
    });

    prevPage.addEventListener('click', () => {
        if (currentPage > 1) {
            currentPage--;
            fetchMedications(searchInput.value, currentPage);
        }
    });

    nextPage.addEventListener('click', () => {
        currentPage++;
        fetchMedications(searchInput.value, currentPage);
    });

    window.editMedication = (id) => {
        const medication = medications.find(m => m.id === id);
        if (medication) {
            modal.style.display = 'block';
            medicationForm.name.value = medication.name;
            medicationForm.quantity.value = medication.quantity;

            medicationForm.onsubmit = (event) => {
                event.preventDefault();
                const updatedMedication = {
                    name: medicationForm.name.value,
                    quantity: medicationForm.quantity.value
                };

                fetch(`http://yourapi.com/medications/${id}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(updatedMedication)
                })
                .then(response => response.json())
                .then(() => {
                    modal.style.display = 'none';
                    fetchMedications(searchInput.value, currentPage);
                })
                .catch(error => {
                    console.error('Error updating medication:', error);  // Логирование ошибок
                });
            };
        }
    };

    window.deleteMedication = (name) => {
        if (confirm('Вы уверены, что хотите удалить этот медикамент?')) {
            fetch('/deleteDrug', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ drug: name })
        })
        .then(response => response.json())
        .then(() => {
            fetchMedications(currentPage, 5);
        })
        .catch(error => {
            console.error('Error deleting medication:', error);
        });
        }
    };

    window.showAdjustModal = (name, action) => {
        adjustMedicationId = name;
        adjustAction = action;
        adjustModal.style.display = 'block';
        adjustForm.reset();
    };

    adjustForm.addEventListener('submit', (event) => {
        event.preventDefault();
        const adjustQuantity = document.getElementById('adjustQuantity').value;
        const url = `http://yourapi.com/medications/${adjustMedicationId}/${adjustAction}`;

        fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ quantity: adjustQuantity })
        })
        .then(response => response.json())
        .then(() => {
            adjustModal.style.display = 'none';
            fetchMedications(searchInput.value, currentPage);
        })
        .catch(error => {
            console.error('Error adjusting medication:', error);  // Логирование ошибок
        });
    });

    fetchMedications();
});