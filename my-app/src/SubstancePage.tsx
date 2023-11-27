import {FC, useEffect, useState} from 'react'
import {Button, Card} from 'react-bootstrap'

import './SubstancePage.css'

import {getSubstanceByName} from './modules/get-substance-by-name.ts'
import {Substance} from './modules/ds'

import defaultImage from './assets/react.svg'

const SubstancePage: FC = () => {

    const [substance, setSubstance] = useState<Substance>()

    useEffect(() => {
        const queryString = window.location.search;
        const urlParams = new URLSearchParams(queryString)
        const SubstanceName = urlParams.get('substance_name')

        const loadRegion = async () => {
            const result: Substance[] = await getSubstanceByName(String(SubstanceName))
            const substanceData = result[0];
           // console.log(result)
            setSubstance(substanceData)
           // console.log(substance)
        }

        loadRegion()

    }, []);

    return (
        <div className='card_container'>
            <Card className='page_card'>
                <Card.Img src={(substance?.Image == '' ? defaultImage?.toString() : "data:image/jpg;base64, " + substance?.Image)} className="card-img-top" variant="top" />                <Card.Body>
                    <p>{substance?.Title && substance.Title}</p>
                    <p> <b>Статус субстанции: {substance?.Status && substance.Status}</b></p>
                    <p> Класс: {substance?.Class && substance.Class}</p>
                    <p> Формула: {substance?.Formula && substance.Formula} </p>
                </Card.Body>
                <Card.Footer>
                    <Button href="/One-pot-front/">Домой</Button>
                </Card.Footer>
            </Card>
        </div>


    )
}

export default SubstancePage