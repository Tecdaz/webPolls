function showMessage(message, isError = false, timeout = 3000) {
  console.log("dentro de show message")
  const container = document.getElementById('messagesContainer');
  if (!container) return;
  console.log("container de mensajes existe")
  const messageDiv = document.createElement('div')
  messageDiv.className = isError ? 'message error' : 'message success';
  messageDiv.textContent = message;
  container.appendChild(messageDiv);

  console.log("mensaje agregado al contenedor")
  requestAnimationFrame(() => {
    messageDiv.classList.add('visible');
  })
  console.log("mensaje visible")
  setTimeout(()=>{
    messageDiv.classList.remove('visible');
    messageDiv.addEventListener('transitionend', () => messageDiv.remove() , { once: true });
  }, timeout);
  console.log("mensaje programado para desaparecer")
}
/*SECCION DE POLLS*/
//renderizado
function renderPolls(polls) {
  const container = document.getElementById('pollsContainer');
  container.innerHTML = '';

  console.log(polls);

  if (!polls || polls.length === 0) {
    container.innerHTML = `<p class="no-polls">No hay encuestas creadas todavía.</p>`;
    return;
  }

  polls.forEach((poll) => {
    const div = document.createElement('div');
    div.classList.add('singlePollDiv');
    div.id = `poll-${poll.poll_id}`;
    
    //revisar esta parte, por que si a veces viene poll.title y otras poll.question esta mal
    const pollId = poll.poll_id || poll.id;
    const title = poll.title || poll.question;
    const options = poll.options || [];

    //seccion de agregado de boton a la opcion
    div.innerHTML = `
      <h3>${title}</h3>
      <ul>
        ${options
          .map(
            (opt) => `
          <li>
            ${opt.content}
            <button class="${opt.correct ? 'markCorrectBtnTrue' : 'markCorrectBtnFalse'}" 
              todavia no se agrego la logica
            </button>
          </li>`
          )
          .join('')}
      </ul>
    `;

    //boton eliminar encuesta
    const deleteBtn = document.createElement('button');
    deleteBtn.textContent = 'Eliminar Encuesta';
    deleteBtn.classList.add('deleteBtn');
    deleteBtn.dataset.id = pollId;
    deleteBtn.addEventListener('click', async () => {
        await deletePoll(pollId);
    });

    div.appendChild(deleteBtn);
    container.appendChild(div);
  });
}

// obtener las encuestas 
async function getPolls() {
  try {
    const res = await fetch('/polls', {
      method: 'GET',
      headers: { 'Accept': 'application/json' },
    });
    if (!res.ok) throw new Error('Error al obtener encuestas');
    
    const data = await res.json();
    const polls = data.data || data; //aca toma la lista de encuestas del campo data
    renderPolls(Array.isArray(polls) ? polls : []);
    
  
  } catch (err) {
    console.error(err);
  } 
}

// creacion de una encuesta
async function createPoll(question, options) {
  try {
    //user id = 3 queda hasta que se aplique las partes de las secciones en los users

    const body = { question, options, user_id: 3 };
   
    const res = await fetch('/polls/create', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json',
        'Accept': 'application/json'
       },
      body: JSON.stringify(body),
    });

    const data = await res.json();
    if (!res.ok) throw new Error(data.message || 'Error al crear encuesta');

     showMessage('Encuesta creada correctamente');
    await getPolls();
  } catch (err) {
    console.error(err);
    showMessage('Error al crear encuesta', true);
  }
}

// eliminar encuesta
async function deletePoll(id) {
  try {
    const res = await fetch(`/polls/${id}`, { method: 'DELETE', headers: { 'Accept': 'application/json' } });
    if (!res.ok) throw new Error('Error al eliminar');

    const pollEl = document.getElementById(`poll-${id}`);
    if (pollEl) {
      pollEl.style.opacity = '0';
      setTimeout(() => {
        pollEl.remove();

        const remaining = document.querySelectorAll('.singlePollDiv').length;
        if (remaining === 0) renderPolls([]);
      }, 300);
    }

    getPolls();
    console.log(pollsData);
   showMessage('Encuesta eliminada correctamente');
  
  } catch (err) {
    console.error(err);
       showMessage('Error al eliminar encuesta', true);
  }
}

// actualizar el estado de una opcion al seleccionar
async function toggleCorrect(optionId, newValue) {
  try {
    const res = await fetch(`/options/${optionId}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json',
        'Accept': 'application/json'
       },
      body: JSON.stringify({ correct: newValue }),
    }); 
    if (!res.ok) throw new Error('Error al actualizar opción')
  } catch (err) {
    console.error('Error en toggleCorrect:', err);
  }
}

// eventos 
document.addEventListener('DOMContentLoaded', () => {
  getPolls();
  getUsers();

  // crear usuario
  const userForm = document.getElementById('userForm');
  userForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = e.target.username.value.trim();
    const email = e.target.email.value.trim();
    const password = e.target.password.value.trim();

    if (!username) return showMessage('formMessage', 'El nombre de usuario no puede estar vacío', true);
    if (!email) return showMessage('formMessage', 'El email no puede estar vacío', true);
    if (!password) return showMessage('formMessage', 'La contraseña no puede estar vacía', true);

    await createUser(username, email, password);
    userForm.reset();
  });

  const form = document.getElementById('pollForm');
  const optsContainer = document.getElementById('optsContainer');
  const addBtn = document.getElementById('addOptBtn');

  // crear
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const question = e.target.question.value.trim();

    const options = Array.from(document.querySelectorAll('input[name="options[]"]'))
      .map((input) => ({ content: input.value.trim(), correct: false }))
      .filter((o) => o.content);

    if (!question) return showMessage('formMessage', 'La pregunta no puede estar vacía', true);
    if (options.length < 2) return showMessage('formMessage', 'Agrega al menos 2 opciones', true);

    await createPoll(question, options);
    form.reset();
    optsContainer.innerHTML = '';
  });

  // agregado de opciones
  addBtn.addEventListener('click', () => {
    if (document.querySelectorAll('input[name="options[]"]').length >= 4) {
      showMessage('formMessage', 'Máximo 4 opciones permitidas', true);
      return;
    }

    const optDiv = document.createElement('div');
    optDiv.classList.add('opt');
    optDiv.innerHTML = `
      <label>Opción</label>
      <input type="text" name="options[]" placeholder="Escribe una opción..." required>
      <button type="button" class="deleteOptBtn">Eliminar</button>
    `;
    optDiv.querySelector('.deleteOptBtn').addEventListener('click', () => optDiv.remove());
    optsContainer.appendChild(optDiv);
  });
});

 // click de botones de seleccion
document.addEventListener('click', async (event) => {
  const btn = event.target;

  // Verificar que sea uno de los botones correctos
  if (!btn.matches('button[data-option-id]')) return;

  const optionId = btn.dataset.optionId;
  const pollId = btn.dataset.pollId;
  const currentState = btn.dataset.correct === 'true';
  const newValue = !currentState;

  // Actualizar visualmente
  btn.classList.remove(currentState ? 'markCorrectBtnTrue' : 'markCorrectBtnFalse');
  btn.classList.add(newValue ? 'markCorrectBtnTrue' : 'markCorrectBtnFalse');
  btn.textContent = newValue ? 'Selected' : 'Select';

  //Actualizar el atributo del dataset (muy importante)
  btn.dataset.correct = String(newValue);

  //Llamar al backend
  try {
    await toggleCorrect(optionId, newValue);
    getPolls()
  } catch (err) {
    console.error('Error al actualizar en servidor:', err);
  }
});

/*SECCION DE USUARIOS*/
async function createUser(username, email, password) {
  try {
    const body = { username, email, password }

    const res = await fetch('/users/create',{
      method: 'POST',
      headers: { 'Content-Type': 'application/json',
        'Accept': 'application/json'
       },
      body: JSON.stringify(body),
    })

    const data = await res.json();
    if (!res.ok) throw new Error(data.message || 'Error al crear usuario');
    getUsers();
    showMessage('Usuario creado correctamente');
  }catch (err) {
    console.error(err);
    showMessage('Error al crear usuario', true);
  }
}

async function getUsers (){
  try {
    const res = await fetch('/users', {
      method: 'GET',
      headers: { 'Accept': 'application/json' },
    });
    if (!res.ok) throw new Error('Error al obtener usuarios');

    const data = await res.json();
    const users = data.data || data; //aca toma la lista de usuarios del campo data
    renderUsers(Array.isArray(users) ? users : []);
  }catch (err) {
    console.error(err);
    showMessage('Error al cargar usuarios', true);
  }
}

async function renderUsers (users){
  const container = document.getElementById('usersContainer');
  container.innerHTML = '';

  if (!users || users.length === 0) {
    container.innerHTML = `<p class="no-users">No hay usuarios registrados todavía.</p>`;
    return;
  }

  users.forEach((user)=>{
    const div = document.createElement('div');
    div.classList.add('singleUserDiv');
    div.id = `user-${user.user_id}`;

    const userId = user.user_id || user.id;
    const username = user.username;
    const email = user.email;

    div.innerHTML = `
      <h3>${username}</h3>
      <p>Email: ${email}</p>
      <p>User ID: ${userId}</p>
    `;

    const deleteBtn = document.createElement('button');
    deleteBtn.textContent = 'Eliminar Encuesta';
    deleteBtn.classList.add('deleteBtn');
    deleteBtn.dataset.id = userId;
    deleteBtn.addEventListener('click', async () => {
        await deleteUser(userId);
    });

    div.appendChild(deleteBtn);
    container.appendChild(div);
  })
}

async function deleteUser (id){

  try {
    const res = await fetch(`/users/${id}`,{
        method: 'DELETE', 
        headers: { 'Accept': 'application/json' } }
    )
    if (!res.ok) throw new Error('Error al eliminar usuario');

    const userEl = document.getElementById(`user-${id}`);
    if (userEl) {
      userEl.style.opacity = '0';
      setTimeout(() => {
        userEl.remove();

        const remaining = document.querySelectorAll('.singleUserDiv').length;
        if (remaining === 0) renderUsers([]);
      }, 300);
    }

    getUsers();
    showMessage('Usuario eliminado correctamente');
  }catch (err) {
    console.error(err);
    showMessage('Error al eliminar usuario', true);
  }
}
