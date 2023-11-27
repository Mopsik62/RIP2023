import {Substance} from './ds'


export const getSubstances = async (namePattern = '') : Promise<Substance[]> => {
    return fetch('/api/substances?name_pattern=' + String(namePattern))
        .then((response) => response.json())
        .catch(() => ([
            {
                "ID": 1,
                "Title": "субстанция1тест",
                "Class": "тест1",
                "Formula": "",
                "Image": "",
                "Status": "Активно"
            },
            {
                "ID": 2,
                "Title": "субстанция2тест",
                "Class": "тест2",
                "Formula": "",
                "Image": "",
                "Status": "Активно"
            },
            {
                "ID": 3,
                "Title": "субстанция3тест",
                "Class": "тест3",
                "Formula": "",
                "Image": "",
                "Status": "Активно"
            },
        ]));
}
