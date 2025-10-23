function showMessage(elementId, message, isError = false) {
  const el = document.getElementById(elementId);
  el.textContent = message;
  el.style.display = 'block';
  el.className = isError ? 'message error' : 'message success';
  setTimeout(() => (el.style.display = 'none'), 3000);
}

//renderizado
function renderPolls(polls) {
  const container = document.getElementById('pollsContainer');
  container.innerHTML = '';

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
              data-option-id="${opt.id}" 
              data-poll-id="${pollId}" 
              data-correct="${opt.correct}">
              ${opt.correct ? 'Selected' : 'Select'}
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
    const res = await fetch('/polls');
    if (!res.ok) throw new Error('Error al obtener encuestas');
    
    const data = await res.json();
    const polls = data.data || data; //aca toma la lista de encuestas del campo data
    renderPolls(Array.isArray(polls) ? polls : []);
  
  } catch (err) {
    console.error(err);
    showMessage('formMessage', 'Error al cargar encuestas', true);
  } 
}

// creacion de una encuesta
async function createPoll(question, options) {
  try {
    //user id = 3 queda hasta que se aplique las partes de las secciones en los users
    const body = { question, options, user_id: 3 };
   
    const res = await fetch('/polls/create', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    });

    const data = await res.json();
    if (!res.ok) throw new Error(data.message || 'Error al crear encuesta');

    showMessage('formMessage', 'Encuesta creada correctamente');
    await getPolls();
  } catch (err) {
    console.error(err);
    showMessage('formMessage', 'Error al crear encuesta', true);
  }
}

// eliminar encuesta
async function deletePoll(id) {
  try {
    const res = await fetch(`/polls/${id}`, { method: 'DELETE' });
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

    showMessage('formMessage', 'Encuesta eliminada correctamente');
  } catch (err) {
    console.error(err);
    showMessage('formMessage', 'Error al eliminar encuesta', true);
  }
}

// actualizar el estado de una opcion al seleccionar
async function toggleCorrect(optionId, newValue) {
  try {
    const res = await fetch(`/options/${optionId}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ correct: newValue }),
    }); 
    if (!res.ok) throw new Error('Error al actualizar opción')
  } catch (err) {
    console.error('Error en toggleCorrect:', err);
  }
}

// eventos globales
document.addEventListener('DOMContentLoaded', () => {
  getPolls();

  const form = document.getElementById('pollForm');
  const optsContainer = document.getElementById('optsContainer');
  const addBtn = document.getElementById('addOptBtn');

  // crear encuesta
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
