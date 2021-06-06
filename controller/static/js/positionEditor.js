
samples = {
    "roll": [],
    "pitch": []
}

async function refreshPosition() {
    const corrected = await getUrl('/api/corrected')
    samples.roll.push(corrected.roll)
    samples.pitch.push(corrected.pitch)
    document.getElementById('rollId').value = median(samples.roll)
    document.getElementById('pitchId').value = median(samples.pitch)
}

async function saveNewPosition() {
    const position = getPositionFromEditor()
    if (validatePosition(position)) {
        const saved = await postUrl('/api/position', position)
        window.location = '/position/' + saved.id
    }
}

async function saveExistingPosition(id) {
    const position = getPositionFromEditor()
    position.id = id
    if (validatePosition(position)) {
        await putUrl('/api/position/' + id, position)
    }
}

async function deletePosition(id) {
    await deleteUrl('/api/position/' + id)
    window.location = '/home'
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
