import tkinter.filedialog as filedialog
import tkinter as tk

master = tk.Tk()
master.title("Primitive")

def input():
    input_path = tk.filedialog.askopenfilename()
    input_entry.delete(1, tk.END)  # Remove current text in entry
    input_entry.insert(0, input_path)  # Insert the 'path'
	
def output():
    path = tk.filedialog.askopenfilename()
    input_entry.delete(1, tk.END)  # Remove current text in entry
    input_entry.insert(0, path)  # Insert the 'path'
	
	
	
	
	
top_frame = tk.Frame(master)
bottom_frame = tk.Frame(master)
line = tk.Frame(master, height=1, width=400, bg="grey80", relief='groove')

	
top_frame.pack(side=tk.TOP)
line.pack(pady=10)
bottom_frame.pack(side=tk.BOTTOM)	
	
	

input_path = tk.Label(top_frame, text="Picture path:")
input_entry = tk.Entry(top_frame, text="", width=40)
browse1 = tk.Button(top_frame, text="Browse", command=input)

output_path = tk.Label(bottom_frame, text="Output Path:")
output_entry = tk.Entry(bottom_frame, text="", width=40)
browse2 = tk.Button(bottom_frame, text="Browse", command=output)
	
begin_button = tk.Button(bottom_frame, text='Begin!') #beginButton	

top_frame.pack(side=tk.TOP)
line.pack(pady=10)
bottom_frame.pack(side=tk.BOTTOM)
	
input_path.pack(pady=5)
input_entry.pack(pady=5)
browse1.pack(pady=5)

output_path.pack(pady=5)
output_entry.pack(pady=5)
browse2.pack(pady=5)

begin_button.pack(pady=20, fill=tk.X)

	
master.mainloop()

#def window():
#       Creates window for the program

#TODO Req 1.0
#def helpButton()
#   Creates a help button

#TODO 1.1
#def export()
#   allow the user to export their finished
#   picture to the location of their choice 
#TODO 1.2
#def difficulty()
#   allows the user to select the number of geometric shapes to form
#   the image

#TODO 1.3
#def geometricShapes()
#   allows the user to select what geometric shape to form
#   the image

#TODO 1.4
#def display()
#   Shows the final form of the image

#TODO Req 1.5
#def file()
    #allows the user to select the photo they want
    #import easygui
    #file = easygui.fileopenbox()
