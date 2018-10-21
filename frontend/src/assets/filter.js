export function phoneNumber (number) {
  return number.substr(0, 3) + '-' + number.substr(3, 4) + '-' + number.substr(7)
}

export function price (number) {
  number = String(Math.ceil(number))
  return number.replace(/(\d)(?=(\d\d\d)+(?!\d))/g, '$1,')
}

export function date (string) {
  const d = new Date(string)
  return d.getFullYear() + '-' + (d.getMonth() + 1) + '-' + d.getDate() + ' ' +
    ('00' + d.getHours()).slice(-2) + ':' +
    ('00' + d.getMinutes()).slice(-2)
}
