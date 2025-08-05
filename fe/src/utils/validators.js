export function validateInput(valor, reglas = {}) {
  if (reglas.minLength && valor.length < reglas.minLength) {
    return true;
  }

  if (reglas.maxLength && valor.length > reglas.maxLength) {
    return true;
  }

  if (reglas.required && valor.length === 0) {
    return true;
  }

  return false;
}

export function removeSpaces(ev, setValor) {
  let valor = ev.target.value;

  valor = valor.replace(/\s/g, "");

  setValor(valor);
  return null;
}
