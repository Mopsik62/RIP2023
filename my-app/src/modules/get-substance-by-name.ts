import {Substance} from './ds'

export const getSubstanceByName = async  (substanceName = ''): Promise<Substance[]> => {
    return fetch('/api/substance/' + String(substanceName),{
        method: 'GET',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        }
    })
        .then((response) => response.json())
        .catch(() => (
            {
                "ID": 1,
                "Title": "субстанция1тест",
                "Class": "тест1",
                "Formula": "",
                "Status": "Активно",
                "Image": ""
            }));
}

