"use client"
import { useEffect, useState } from "react";
import { ToastContainer, toast } from "react-toastify";

export default function Home() {
  const [user, setUser] = useState('')
  const [creci, setCreci] = useState('')
  const [content, setContent] = useState('')
  const [parecerId, setParecerId] = useState('')

  const [editParecer, setEditParecer] = useState(false)
  const [pareceres, setPareceres] = useState([])
  const [reload, setReload] = useState(false)

  const notify = (error) => {
    toast(error)

  }

  useEffect(() => {
    fetch('https://parecer-gen.onrender.com/parecer')
      .then(res => res.json())
      .then(data => setPareceres(data))
      .catch(err => notify('Erro ao carregar pareceres'))
  }, [reload])


  const handleReload = () => setReload(!reload);

  return (
    <>
      <ToastContainer aria-label={undefined} />

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
          {!editParecer ? <button onClick={() => {
            fetch('https://parecer-gen.onrender.com/parecer', {
              method: 'POST',
              body: JSON.stringify({
                user, creci, content
              })
            }).
              then(res => {
                if (res.status != 201) {
                  notify("Erro no para criar parecer, codigo " + res.status)
                }
                setUser('')
                setCreci('')
                setContent('')
                setParecerId('')
                handleReload()
              }).
              catch(err => notify('Erro ao mandar parecer ' + err))
          }} className="flex border rounded-xl dark:bg-gray-700 bg-gray-100 py-2 px-4">Gerar Parecer</button> :
            <button onClick={() => {
              fetch(`https://parecer-gen.onrender.com/parecer?id=${parecerId}`, {
                method: 'PUT',
                body: JSON.stringify({
                  user, creci, content
                })
              }).then(res => {
                if (res.status != 200) {
                  notify("Erro no servidor para atualizar parecer")
                }
                setUser('')
                setCreci('')
                setContent('')
                setParecerId('')
                handleReload()
              }).catch(err => {
                notify('Erro ao editar parecer ' + err)
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
            <li key={parecer.id} className="flex gap-4  w-full">
              <div className="p-4 border rounded bg-gray-100 w-full dark:bg-gray-700">
                <p><strong>Usuário:</strong> {parecer.user}</p>
                <p><strong>Creci:</strong> {parecer.creci}</p>
                <p><strong>Data:</strong> {parecer.date}</p>
              </div>
              <div className="flex flex-col md:flex-row gap-4 justify-evenly">
                <a className="flex items-center p-4 border rounded bg-gray-100 dark:bg-gray-700 " href={`https://parecer-gen.onrender.com/parecer?id=${parecer.id}`} target="_blank">Baixar</a>
                <div className="flex items-center p-4 border rounded bg-gray-100 dark:bg-gray-700" onClick={() => {
                  setUser(parecer.user)
                  setCreci(parecer.creci)
                  setContent(parecer.content)
                  setParecerId(parecer.id)
                  setEditParecer(true)
                }}>Editar</div>
                <div className="flex items-center p-4 border rounded bg-gray-100 dark:bg-gray-700" onClick={() => {
                  fetch(`https://parecer-gen.onrender.com/parecer?id=${parecer.id}`, {
                    method: 'DELETE',
                  }).then(res => {
                    if (res.status != 200) {
                      notify("Erro no servidor para deletar parecer")
                    }
                    handleReload()
                  }).catch(err => notify('Erro ao deletar parecer ' + err))
                }}>Deletar</div>
              </div>
            </li>
          ))}
        </ul>
      </div >

    </>
  );
}

