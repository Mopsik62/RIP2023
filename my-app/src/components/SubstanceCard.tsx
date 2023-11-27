import {FC} from 'react'
import {Button, Card} from 'react-bootstrap'
import './SubstanceCard.css'

interface Props {
    imageUrl: string
    SubstanceName: string
    pageUrl: string
}



const SubstanceCard: FC<Props> = ({ imageUrl, SubstanceName, pageUrl}) => {

    const deleteRestoreRegion = async () => {
        await fetch('/api/substance/delete/' + SubstanceName, {
            method: 'PUT'
        });
        window.location.replace('/')
    }

    return (
        <Card>
            <Card.Img className="card-img-top" variant="top" src={imageUrl}/>
            <Card.Body>
                <div className='textStyle'>
                    <Card.Title> {SubstanceName} </Card.Title>
                </div>
            </Card.Body>
            <Card.Footer>
                <div className="btn-wrapper text-center d-flex justify-content-between">
                    <Button variant="secondary" href={pageUrl}>Подробнее</Button>
                    <Button variant="warning" onClick={deleteRestoreRegion}>Изменить статус</Button>
                </div>
            </Card.Footer>
        </Card>

    )

}

export default SubstanceCard;