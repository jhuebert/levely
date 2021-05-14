
function displayEditor(position) {

    if (isNewPosition(position)) {
        position.name = new Date().toISOString()
    }

    return `
        <div class="container h-100">
            <div class="row px-3" id="alert"></div>
            <div class="text-light text-center h3 mb-2">${isNewPosition(position) ? 'New' : 'Edit'} Position</div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Name</label>
                <div class="col-sm-10">
                    <input type="text" class="form-control" id="nameId" placeholder="Enter position name" value="${position.name}">
                </div>
            </div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Roll</label>
                <div class="col-sm-10">
                    <input type="number" class="form-control" id="rollId" placeholder="0.0" value="${position.roll}">
                </div>
            </div>
            <div class="form-group row align-items-center">
                <label class="col-sm-2 col-form-label col-form-label-lg text-light">Pitch</label>
                <div class="col-sm-10">
                    <input type="number" class="form-control" id="pitchId" placeholder="0.0" value="${position.pitch}">
                </div>
            </div>
            <div class="custom-control custom-switch">
              <input type="checkbox" class="custom-control-input" id="favoriteId" ${position.favorite ? 'checked' : ''}>
              <label class="custom-control-label text-light h5" for="favoriteId">Favorite</label>
            </div>
            <button type="button" class="btn btn-primary btn-block mt-5" onclick="${isNewPosition(position) ? `saveNewPosition()` : `saveExistingPosition(${position.id})`}">Save</button>
            <button type="button" class="btn btn-danger btn-block" data-toggle="modal" data-target="#deleteModal" ${isNewPosition(position) ? 'hidden' : ''}>Delete</button>
            <div class="modal fade" id="deleteModal" tabindex="-1" role="dialog" aria-labelledby="deleteModalLabel" aria-hidden="true">
              <div class="modal-dialog" role="document">
                <div class="modal-content">
                  <div class="modal-header">
                    <h5 class="modal-title" id="deleteModalLabel">Delete Position</h5>
                    <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                      <span aria-hidden="true">&times;</span>
                    </button>
                  </div>
                  <div class="modal-body">Are you sure you want to delete this position?</div>
                  <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
                    <button type="button" class="btn btn-danger" onclick="deletePosition(${position.id})" data-dismiss="modal">Delete</button>
                  </div>
                </div>
              </div>
            </div>
        </div>
    `;
}

function isNewPosition(position) {
    return !position.id || position.id === 0
}

async function displaySaveNew(matches, div) {
    const position = await getUrl('api/corrected')
    div.innerHTML = displayEditor(position);
}

async function displayEditExiting(matches, div) {
    const position = await getUrl('api/position/' + matches[1])
    div.innerHTML = displayEditor(position);
}

async function saveNewPosition() {
    const position = getPositionFromEditor()
    if (validatePosition(position)) {
        const saved = await postUrl('api/position', position)
        window.location = '#position/' + saved.id
    }
}

async function saveExistingPosition(id) {
    const position = getPositionFromEditor()
    position.id = id
    if (validatePosition(position)) {
        await putUrl('api/position/' + id, position)
        displayAlert('alert-success', 'Update success', true)
    }
}

async function deletePosition(id) {
    await deleteUrl('api/position/' + id)
    window.location = '#home'
}

function getPositionFromEditor() {
    const name = document.getElementById('nameId').value;
    const roll = parseFloat(document.getElementById('rollId').value);
    const pitch = parseFloat(document.getElementById('pitchId').value);
    const favorite = document.getElementById('favoriteId').checked;
    return {
        name: name,
        roll: roll,
        pitch: pitch,
        favorite: favorite
    }
}

function validatePosition(position) {
    if (position.name === '') {
        displayAlert('alert-danger', 'Name is required', false)
        return false
    }
    return true
}
