function showMessage(message, isError = false, timeout = 3000) {
  console.log("dentro de show message")
  const container = document.getElementById('messagesContainer');
  if (!container) return;
  console.log("container de mensajes existe")

  //div dinamico para contener el texto del mensaje
  const messageDiv = document.createElement('div')
  
  //asignamos clases dependiendo si es error o exito para los estilos
  messageDiv.className = isError ? 'message error' : 'message success';
  messageDiv.textContent = message;
  container.appendChild(messageDiv); //insertamos en DOM para hacer visible el mensaje en la pagina

  console.log("mensaje agregado al contenedor")
  requestAnimationFrame(() => {
    messageDiv.classList.add('visible');
  })
  console.log("mensaje visible")
  setTimeout(()=>{ //animaciones de entrada y salida
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
    //algunos backends pueden devolver poll_id o id, por eso se contemplan ambos
    const pollId = poll.poll_id || poll.id;
    const title = poll.title || poll.question;
    const options = poll.options || []; //si no hay opciones se evita un error

    //seccion de agregado de boton para maracar la opcion -> pendiente
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
    //cada una tiene su propio listener
    const deleteBtn = document.createElement('button');
    deleteBtn.textContent = 'Eliminar Encuesta';
    deleteBtn.classList.add('deleteBtn');
    deleteBtn.dataset.id = pollId;
    deleteBtn.addEventListener('click', async () => {
        await deletePoll(pollId);
    });

    div.appendChild(deleteBtn);
    container.appendChild(div); //agrega encuesta completa al DOM
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
    const polls = data.data; //saque || data porque el backend siempre devuelve envuelto en data
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
    console.log(data);
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

    //animacion de salida antes de eliminar del DOM
    if (pollEl) {
      pollEl.style.opacity = '0';
      setTimeout(() => {
        pollEl.remove();

        const remaining = document.querySelectorAll('.singlePollDiv').length;
        if (remaining === 0) renderPolls([]);
      }, 300);
    }

    getPolls(); //una vez eliminada se recargan todas las encuestas
    showMessage('Encuesta eliminada correctamente');
  
  } catch (err) {
    console.error(err);
       showMessage('Error al eliminar encuesta', true);
  }
}

// actualizar el estado de una opcion al seleccionar -> pendiente, se puede comentar o sacar
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
  const pollsContainer = document.getElementById('pollsContainer');
  const usersContainer = document.getElementById('usersContainer');
  
  //cargamos los datos dependiendo la pagina donde esta el usuario
  if (pollsContainer) getPolls();
  if (usersContainer) getUsers();

  // crear usuario
  const userForm = document.getElementById('userForm');
  if (userForm){
    userForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const username = e.target.username.value.trim();
      const email = e.target.email.value.trim();
      const password = e.target.password.value.trim();

      if (!username){
        showMessage('El nombre de usuario no puede estar vacío', true);
        return;
      }
      if (!email){
        showMessage('El email no puede estar vacío', true);
        return;
      }
      if (!password){
        showMessage('La contraseña no puede estar vacía', true);
        return;
      }

      try {
        await createUser(username, email, password);
        userForm.reset(); //limpiar despues del envio
      } catch (err) {
        console.error(err);
        showMessage('Error al crear usuario', true);
      }
    });
  }

  const pollForm = document.getElementById('pollForm');
  const optsContainer = document.getElementById('optsContainer');
  const addBtn = document.getElementById('addOptBtn');

  if (pollForm && optsContainer && addBtn){
    // crear
    // los inputs dinamicos de opciones se convierten en arreglos de opciones para poder usar map y filter
    pollForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const question = e.target.question.value.trim();

      const options = Array.from(document.querySelectorAll('input[name="options[]"]'))
        .map((input) => ({ content: input.value.trim(), correct: false })) //convertir en objeto porque es lo que espera la API
        .filter((o) => o.content); //filtrar opciones vacias

      if (!question){
        showMessage('La pregunta no puede estar vacía', true);
        return;
      }
      if (options.length < 2){
        console.log(options);
        showMessage('Agrega al menos 2 opciones', true);
        return;
      }

      try {
        await createPoll(question, options);
        pollForm.reset();
        optsContainer.innerHTML = '';
      } catch (err) {
        console.error(err);
        showMessage('Error al crear encuesta', true);
      }
    });

    // agregado de opciones
    addBtn.addEventListener('click', () => {
      if (document.querySelectorAll('input[name="options[]"]').length >= 4) {
        showMessage('Máximo 4 opciones permitidas', true);
        return;
      }

      const optDiv = document.createElement('div'); //contenedor de cada opcion
      optDiv.classList.add('opt');
      optDiv.innerHTML = `
        <label>Opción</label>
        <input type="text" name="options[]" placeholder="Escribe una opción..." required>
        <button type="button" class="deleteOptBtn">Eliminar</button>
      `;
      optDiv.querySelector('.deleteOptBtn').addEventListener('click', () => optDiv.remove()); //agregamos el evento al boton recien creado
      optsContainer.appendChild(optDiv);
    });
  }
});

// todo esto pertenece a la logica del boton de seleccion que todavia no implementamos -> pendiente
 // click de botones de seleccion
//   document.addEventListener('click', async (event) => {
//   const btn = event.target;

//   //verificar que sea uno de los botones correctos
//   if (!btn.matches('button[data-option-id]')) return;

//   const optionId = btn.dataset.optionId;
//   const pollId = btn.dataset.pollId;
//   const currentState = btn.dataset.correct === 'true';
//   const newValue = !currentState;

//   //actualizar visualmente
//   btn.classList.remove(currentState ? 'markCorrectBtnTrue' : 'markCorrectBtnFalse');
//   btn.classList.add(newValue ? 'markCorrectBtnTrue' : 'markCorrectBtnFalse');
//   btn.textContent = newValue ? 'Selected' : 'Select';

//   //actualizar el atributo del dataset (muy importante)
//   btn.dataset.correct = String(newValue);

//   //llamar al backend
//   try {
//     await toggleCorrect(optionId, newValue);
//     getPolls()
//   } catch (err) {
//     console.error('Error al actualizar en servidor:', err);
//   }
// });

// termina la logica del boton de seleccion

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
    getUsers(); //recarga usuarios para mostrar el nuevo
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
    const users = data.data; //borre || data de aca tambien
    renderUsers(Array.isArray(users) ? users : []);
  }catch (err) {
    console.error(err);
    showMessage('Error al cargar usuarios', true);
  }
}

async function renderUsers (users){
  const container = document.getElementById('usersContainer');
  container.innerHTML = '';

  //Eliminar el user con id = 3
  users = users.filter(user => user.id   !== 3);
  console.log(users);

  if (!users || users.length === 0) {
    container.innerHTML = `<p class="no-users">No hay usuarios registrados todavía.</p>`;
    return;
  }

  //creacion de contenedor para cada usuario
  users.forEach((user)=>{
    const div = document.createElement('div');
    div.classList.add('singleUserDiv');
    div.id = `user-${user.id}`;

    const userId = user.id; //porque en user service y pollservice tenemos diferentes convenciones
    const username = user.username;
    const email = user.email;
    //insrtar la info del usuario en el div
    div.innerHTML = ` 
      <h3>${username}</h3>
      <p>Email: ${email}</p>
      <p>User ID: ${userId}</p>
    `;

    //boton eliminar para el usuario
    const deleteBtn = document.createElement('button');
    deleteBtn.textContent = 'Eliminar Usuario';
    deleteBtn.classList.add('deleteBtn');
    deleteBtn.dataset.id = userId;
    deleteBtn.addEventListener('click', async () => { //listener para el evento click que llama a deletePoll
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
        if (remaining === 0) renderUsers([]); //si era el ultimo mostramos que no hay usuarios
      }, 300);
    }

    getUsers(); //se recargan todos los usuarios del servidor
    showMessage('Usuario eliminado correctamente');
  }catch (err) {
    console.error(err);
    showMessage('Error al eliminar usuario', true);
  }
}
