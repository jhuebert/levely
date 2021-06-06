
samples = {
    "roll": [],
    "pitch": []
}

async function saveCalibration() {
    const roll = parseFloat(document.getElementById('rollId').value)
    const pitch = parseFloat(document.getElementById('pitchId').value)
    await putUrl('/api/calibration', { roll: roll, pitch: pitch })
}

async function refreshCalibration() {
    const uncorrected = await getUrl('/api/uncorrected')
    samples.roll.push(uncorrected.roll)
    samples.pitch.push(uncorrected.pitch)
    document.getElementById('rollId').value = median(samples.roll)
    document.getElementById('pitchId').value = median(samples.pitch)
}
