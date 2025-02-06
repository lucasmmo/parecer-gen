"use client"
import { useEffect, useState } from "react";
import { ToastContainer, toast } from "react-toastify";

export default function Home() {
  const [error, setError] = useState('')

  const [user, setUser] = useState('')
  const [creci, setCreci] = useState('')
  const [content, setContent] = useState('')
  const [parecerId, setParecerId] = useState('')

  const [editParecer, setEditParecer] = useState(false)
  const [pareceres, setPareceres] = useState([])
  const [created, setCreated] = useState(false)

  const notify = (error) => {
    toast(error)

  }

  useEffect(() => {
    fetch('http://localhost:8080/parecer')
      .then(res => res.json())
      .then(data => setPareceres(data))
      .catch(err => notify('Erro ao carregar pareceres'))
  }, [created])


  return (
    <>
      <ToastContainer />

      <h1 className="text-3xl">Gerador de Parecer</h1>

      <div className="flex flex-col justify-center w-full px-10 sm:w-1/2 gap-4">
        <div className="flex flex-col">
          <label>Usuário:</label>
          <input className="rounded p-2 dark:bg-gray-700 bg-gray-100" type="text" name="user" id="user" value={user} onChange={(e) => setUser(e.target.value)} />
        </div>

        <div className="flex flex-col">
          <label>Creci:</label>
          <input className="rounded p-2 dark:bg-gray-700 bg-gray-100" type="text" name="creci" id="creci" value={creci} onChange={(e) => setCreci(e.target.value)} />
        </div>

        <div className="flex flex-col">
          <label>Conteúdo:</label>
          <textarea className="rounded p-2 dark:bg-gray-700 bg-gray-100" name="content" id="content" rows={10} value={content} onChange={(e) => setContent(e.target.value)}></textarea>
        </div>

        <div className="flex justify-center items-center my-4">
          {!editParecer ? <button onClick={async () => {
            fetch('http://localhost:8080/parecer', {
              method: 'POST',
              body: JSON.stringify({
                user, creci, content
              })
            }).
              then(res => {
                if (res.status != 201) {
                  notify("Erro no servidor para criar parecer")
                }
                setUser('')
                setCreci('')
                setContent('')
                setParecerId('')
                setCreated(!created)
              }).
              catch(err => notify('Erro ao gerar parecer'))
          }} className="flex border rounded-xl dark:bg-gray-700 bg-gray-100 py-2 px-4">Gerar Parecer</button> :
            <button onClick={async () => {
              fetch(`http://localhost:8080/parecer?id=${parecerId}`, {
                method: 'PUT',
                body: JSON.stringify({
                  user, creci, content
                })
              }).catch(err => {
                notify('Erro ao editar parecer')
                setUser('')
                setCreci('')
                setContent('')
                setParecerId('')
                setEditParecer(false)
              })
            }} className="flex border rounded-xl dark:bg-gray-700 bg-gray-100 py-2 px-4">Editar Parecer</button>
          }
        </div>

        <ul className="flex justify-center flex-col gap-y-4">
          {pareceres && pareceres.map((parecer) => (
            <li key={parecer.id} className="flex justify-between gap-4">
              <div className="p-4 border rounded bg-gray-100 w-full dark:bg-gray-700">
                <p><strong>Usuário:</strong> {parecer.user}</p>
                <p><strong>Creci:</strong> {parecer.creci}</p>
                <p><strong>Data:</strong> {parecer.date}</p>
              </div>
              <div className="flex flex-col md:flex-row gap-4 justify-center">
                <a className="flex items-center p-4 border rounded bg-gray-100 dark:bg-gray-700 " href={`http://localhost:8080/parecer?id=${parecer.id}`} target="_blank">Baixar</a>
                <div className="flex items-center p-4 border rounded bg-gray-100 dark:bg-gray-700" onClick={() => {
                  setUser(parecer.user)
                  setCreci(parecer.creci)
                  setContent(parecer.content)
                  setParecerId(parecer.id)
                  setEditParecer(true)
                }}>Editar</div>
              </div>
            </li>
          ))}
        </ul>
      </div >

    </>
  );
}

