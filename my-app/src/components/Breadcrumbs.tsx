import Breadcrumb from 'react-bootstrap/Breadcrumb'

import './Breadcrumbs.css'

function Breadcrumbs() {

    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString)
    const substance_name = urlParams.get('substance_name')
    const name_pattern = urlParams.get('name_pattern')

    return (
        <Breadcrumb>
            <Breadcrumb.Item href="/One-pot-front/">Домашняя страница</Breadcrumb.Item>
            {(substance_name != null && name_pattern === null) &&
                <>
                    <Breadcrumb.Item active> Субстанция </Breadcrumb.Item>
                    <Breadcrumb.Item href = {window.location.search}>{substance_name}</Breadcrumb.Item>
                </>
            }
            {(name_pattern != null && substance_name === null) &&
                <>
                    <Breadcrumb.Item active> Поиск </Breadcrumb.Item>
                    <Breadcrumb.Item href = {window.location.search}>{name_pattern}</Breadcrumb.Item>
                </>

            }
        </Breadcrumb>
    );
}

export default Breadcrumbs;